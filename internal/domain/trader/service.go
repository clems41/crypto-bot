package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/internal/service/tradingStrategy"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/tradingUtils"
	"fmt"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	Start() (err error)
	Stop() (err error)

	/* Private methods */
	applyTradingAlgorithm() (err error)
	updateBalance(platform tradingPlatform.Api) (err error)
	updateIndexPrice(platform tradingPlatform.Api) (err error)
	addOrder(platform tradingPlatform.Api, order *model.Order) (err error)
	fillOrder(platform tradingPlatform.Api, order *model.Order) (err error)
}

type service struct {
	quitChannel                        chan bool                          // quit goroutine when program exit
	config                             *Config                            // config that should be applied with trading algorithm
	platformApis                       map[string]tradingPlatform.Api     // communicate with trading platforms
	repo                               repository.Repository              // use to store data
	algo                               tradingStrategy.Algo               // use to know if order should be open based on prices
	startTime                          time.Time                          // datetime when lago has been started
	initialBalanceByPlatformByCurrency map[string]map[string]float64      // initial balance before opening first position by platform and by currency
	balanceByPlatform                  map[string]*model.Balance          // current balance
	indexPriceByPlatformByPair         map[string]map[string]*model.Price // store last index got from platform
	balanceNeedToBeUpdated             bool                               // use to know if request should be sent to platform to know actual balance
}

func NewService(config *Config, platformApis []tradingPlatform.Api, repo repository.Repository, algo tradingStrategy.Algo) (Service, error) {
	if config == nil {
		defaultConfig, err := GetConfigFromEnvOrDefault()
		if err != nil {
			return nil, err
		}
		config = &defaultConfig
	}
	platformApisMap := make(map[string]tradingPlatform.Api)
	for _, platformApi := range platformApis {
		platformApisMap[platformApi.Name()] = platformApi
	}
	return &service{
		config:                             config,
		quitChannel:                        make(chan bool),
		platformApis:                       platformApisMap,
		repo:                               repo,
		algo:                               algo,
		initialBalanceByPlatformByCurrency: make(map[string]map[string]float64),
		balanceByPlatform:                  make(map[string]*model.Balance),
		indexPriceByPlatformByPair:         make(map[string]map[string]*model.Price),
		balanceNeedToBeUpdated:             true,
	}, nil
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()

	// init balance
	for platformName, _ := range svc.platformApis {
		if svc.balanceByPlatform[platformName] == nil {
			svc.balanceByPlatform[platformName] = &model.Balance{
				PlatformName:    platformName,
				ValueByCurrency: make(map[string]float64),
			}
		}
	}

	// Set initial balance, useful to calculate ending profit, result, etc...
	for platformName, platform := range svc.platformApis {
		err = platform.RefreshBalance(svc.balanceByPlatform[platformName])
		if err != nil {
			return
		}
		for currency, value := range svc.balanceByPlatform[platformName].ValueByCurrency {
			if svc.initialBalanceByPlatformByCurrency[platform.Name()] == nil {
				svc.initialBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
			}
			svc.initialBalanceByPlatformByCurrency[platform.Name()][currency] = value
		}
	}

	// Run algorithm each X ms
	for range time.Tick(time.Duration(svc.config.DelayBetweenEachRunInMilliSeconds) * time.Millisecond) { // Loop
		select {
		case <-svc.quitChannel:
			logger.Debugf("Stop message has been received")
			return
		default:
			err = svc.applyTradingAlgorithm()
			if err != nil {
				logger.Error(err)
				return
			}
		}
	}
	return
}

func (svc *service) Stop() (err error) {
	svc.quitChannel <- true

	// Cancel all open orders
	for platformName, platform := range svc.platformApis {
		logger.Infof("--  %s  --", platformName)
		var view tradingPlatform.CancelAllOrdersView
		view, err = platform.CancelAllOrders()
		if err != nil {
			return
		}
		logger.Infof("Closing %d orders", view.Count)

		// Get all orders
		var orders []model.Order
		orders, err = platform.GetAllOrders()
		if err != nil {
			return
		}
		for _, order := range orders {
			logger.Infof("Order %v", order)
		}

		// Get current balance
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}

		// Calculate estimated profit
		endTime := time.Now()
		tradingDuration := endTime.Sub(svc.startTime)
		var initialBalance, finalBalance float64
		for currency, initialBalanceCurrency := range svc.initialBalanceByPlatformByCurrency[platform.Name()] {
			initialBalance += initialBalanceCurrency
			finalBalanceCurrency, ok := svc.balanceByPlatform[platformName].ValueByCurrency[currency]
			if !ok {
				return fmt.Errorf("cannot find balance for currency %s", currency)
			}
			finalBalance += finalBalanceCurrency
		}
		oneDayProfit := tradingUtils.EstimateProfit(initialBalance, finalBalance, tradingDuration, 24*time.Hour)
		oneMonthProfit := tradingUtils.EstimateProfit(initialBalance, finalBalance, tradingDuration, 30*24*time.Hour)
		oneYearProfit := tradingUtils.EstimateProfit(initialBalance, finalBalance, tradingDuration, 365*24*time.Hour)
		logger.Infof("With initial balance of %0.2f, profit for one day would be %0.2f, for one month %0.2f and for one year %0.2f",
			initialBalance, oneDayProfit, oneMonthProfit, oneYearProfit)
	}
	return
}

func (svc *service) applyTradingAlgorithm() (err error) {
	logger.Infof("-------------------  New run %s  ------------------------------", time.Now().Format(time.RFC3339))
	for platformName, platform := range svc.platformApis {
		logger.Infof("--  %s  --", platformName)

		// update prices for all pairs and platforms
		err = svc.updateIndexPrice(platform)
		if err != nil {
			return
		}

		// update balance
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}

		// Loop over all pairs to trade
		for _, pair := range pairsToTradeByPlatform[platformName] {
			// get prices history based on interval config
			var prices []model.Price
			priceForm := repository.GetPriceHistoryForm{
				PlatformName: platformName,
				Pair:         pair,
				SinceTime:    time.Now().Add(-time.Duration(svc.config.IntervalToComparePricesInMinutes) * time.Minute),
			}
			prices, err = svc.repo.GetPriceHistory(priceForm)
			if err != nil {
				return
			}

			// If price history doesn't return enough prices, svc.algo.ShouldAddOrder will return an error.
			// So we skip this run, and try again the next one.
			// It can take some time to get enough prices at start, depending on svc.config.IntervalToComparePricesInMinutes.
			if len(prices) < svc.algo.PricesNeeded() {
				logger.Debugf("Run for pair %s will be skip, doesn't get enough prices from history to know if order should be open", pair)
				continue
			}

			// fill open form to know if new order should be open
			indexPrice, ok := svc.indexPriceByPlatformByPair[platformName][pair]
			if !ok || indexPrice == nil {
				return fmt.Errorf("cannot get index price for platform %s and pair %s", platformName, pair)
			}
			openForm := tradingStrategy.ShouldAddOrderForm{
				PriceHistory: prices,
				IndexPrice:   *indexPrice,
			}
			var openView tradingStrategy.ShouldAddOrderView
			openView, err = svc.algo.ShouldAddOrder(openForm)
			if err != nil {
				return
			}

			// open order if conditions are ok
			if openView.ShouldAddOrder {
				order := model.Order{
					Pair:  pair,
					Side:  openView.Side,
					Type:  openView.Type,
					Price: openView.Price,
				}
				err = svc.addOrder(platform, &order)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (svc *service) updateBalance(platform tradingPlatform.Api) (err error) {
	if !svc.balanceNeedToBeUpdated {
		return
	}
	err = platform.RefreshBalance(svc.balanceByPlatform[platform.Name()])
	if err != nil {
		return
	}
	err = svc.repo.StoreBalance(svc.balanceByPlatform[platform.Name()])
	if err != nil {
		return
	}

	// init balance if not already done
	if svc.initialBalanceByPlatformByCurrency[platform.Name()] == nil {
		svc.initialBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
		svc.initialBalanceByPlatformByCurrency[platform.Name()] = svc.balanceByPlatform[platform.Name()].ValueByCurrency
	}
	logger.Infof("Balance %v at %s", svc.balanceByPlatform[platform.Name()].ValueByCurrency,
		svc.balanceByPlatform[platform.Name()].UpdatedAt.Format(time.RFC3339))
	svc.balanceNeedToBeUpdated = false
	return
}

func (svc *service) updateIndexPrice(platform tradingPlatform.Api) (err error) {
	pairs, ok := pairsToTradeByPlatform[platform.Name()]
	if !ok {
		return fmt.Errorf("cannot find pairs for platform %s", platform.Name())
	}
	prices, err := platform.GetIndexPrices(pairs...)
	if err != nil {
		return
	}
	for _, price := range prices {
		err = svc.repo.StorePrice(&price)
		if err != nil {
			return
		}
		if svc.indexPriceByPlatformByPair[platform.Name()] == nil {
			svc.indexPriceByPlatformByPair[platform.Name()] = make(map[string]*model.Price)
		}
		svc.indexPriceByPlatformByPair[platform.Name()][price.Pair] = &price
		logger.Infof("%s ask=%f bid=%f at %s", price.Pair, price.Ask,
			price.Bid, price.Date.Format(time.RFC3339))
	}
	return
}

func (svc *service) addOrder(platform tradingPlatform.Api, order *model.Order) (err error) {
	// fill missing order fields
	err = svc.fillOrder(platform, order)
	if err != nil {
		return
	}

	// don't open order if balance is less or equal to 0
	if order.Amount <= 0 {
		return
	}

	// add order using platform
	err = platform.AddOrder(order)
	if err != nil {
		return
	}
	svc.balanceNeedToBeUpdated = true

	// store order into repository
	err = svc.repo.StoreOrder(order)
	if err != nil {
		return
	}
	logger.Infof("New order has been added %v", *order)
	return
}

func (svc *service) fillOrder(platform tradingPlatform.Api, order *model.Order) (err error) {
	// get current balance needed for pair to trade
	currency, ok := tradingUtils.CurrencyNeededToTradePair(order.Pair, order.Side)
	if !ok {
		return fmt.Errorf("cannot find currency for pair %s and side %s", order.Pair, order.Side)
	}
	balance, ok := svc.balanceByPlatform[platform.Name()].ValueByCurrency[currency]
	if !ok {
		return fmt.Errorf("cannot find balance for currency %s", currency)
	}

	// fill volume based on price and amount
	if order.Side == tradingConst.BuySideOrder {
		order.Amount = balance
		order.Volume = balance / order.Price
	} else {
		order.Amount = balance * order.Price
		order.Volume = balance
	}

	return
}
