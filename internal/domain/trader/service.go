package trader

import (
	"crypto-bot/external/service/tradingPlatform"
	"crypto-bot/internal/constant/timeConst"
	"crypto-bot/internal/constant/tradingConst"
	tradingStrategy2 "crypto-bot/internal/domain/tradingStrategy"
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
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
	getTradeInfo(platform tradingPlatform.Api) (info TradeInfo, err error)
	updateOpenedOrders(platform tradingPlatform.Api) (err error)
	updateBalance(platform tradingPlatform.Api) (err error)
	updateIndexPrice(platform tradingPlatform.Api, pairs []string) (err error)
	addOrder(platform tradingPlatform.Api, order *model.Order) (err error)
}

type service struct {
	quitChannel                        chan bool                           // quit goroutine when program exit
	platformApis                       map[string]tradingPlatform.Api      // communicate with trading platforms
	repo                               repository.Repository               // use to store data
	algo                               tradingStrategy2.Algo               // use to know if order should be open based on prices
	startTime                          time.Time                           // datetime when lago has been started
	initialBalanceByPlatformByCurrency map[string]map[string]float64       // initial balance before opening first order by platform and by currency
	balanceByPlatform                  map[string]model.Balance            // current balance
	indexPriceByPlatformByPair         map[string]map[string]model.Price   // store last index got from platform
	openedOrdersByPlatformByPair       map[string]map[string][]model.Order // store opened orders calculating amount to invest by pair
}

func NewService(platformApis []tradingPlatform.Api, repo repository.Repository, algo tradingStrategy2.Algo) (Service, error) {
	platformApisMap := make(map[string]tradingPlatform.Api)
	for _, platformApi := range platformApis {
		platformApisMap[platformApi.Name()] = platformApi
	}
	return &service{
		quitChannel:                        make(chan bool),
		platformApis:                       platformApisMap,
		repo:                               repo,
		algo:                               algo,
		initialBalanceByPlatformByCurrency: make(map[string]map[string]float64),
		balanceByPlatform:                  make(map[string]model.Balance),
		indexPriceByPlatformByPair:         make(map[string]map[string]model.Price),
		openedOrdersByPlatformByPair:       make(map[string]map[string][]model.Order),
	}, nil
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()

	// init balance
	for platformName, _ := range svc.platformApis {
		svc.balanceByPlatform[platformName] = model.Balance{
			PlatformName:    platformName,
			ValueByCurrency: make(map[string]float64),
		}
	}

	// Set initial balance, useful to calculate ending profit, result, etc...
	for platformName, platform := range svc.platformApis {
		err = svc.updateBalance(platform)
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
	for range time.Tick(delayBetweenEachRun) { // Loop
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
		// Get all orders
		var orders []model.Order
		orders, err = platform.GetAllOrders()
		if err != nil {
			return
		}
		for _, order := range orders {
			logger.Info(order)
		}

		// Get current balance
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}

		// Calculate estimated profit
		endTime := time.Now()
		tradingDuration := endTime.Sub(svc.startTime)
		for currency, initialBalanceCurrency := range svc.initialBalanceByPlatformByCurrency[platform.Name()] {
			finalBalanceCurrency, ok := svc.balanceByPlatform[platformName].ValueByCurrency[currency]
			if !ok {
				return fmt.Errorf("cannot find balance for currency %s", currency)
			}
			oneDayProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 24*time.Hour)
			oneMonthProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 30*24*time.Hour)
			oneYearProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 365*24*time.Hour)
			logger.Infof("With initial balance for currency %s of %f, profit for one day would be %f, for one month %f and for one year %f",
				currency, initialBalanceCurrency, oneDayProfit, oneMonthProfit, oneYearProfit)
		}
	}
	return
}

func (svc *service) applyTradingAlgorithm() (err error) {
	logger.Infof("-------------------  New run %s  ------------------------------", time.Now().Format(timeConst.DefaultFormatTimeLayout))
	for platformName, platform := range svc.platformApis {
		logger.Infof("--  %s  --", platformName)
		// update balance
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}

		// update open orders form platform
		err = svc.updateOpenedOrders(platform)
		if err != nil {
			return
		}

		// get trade info (pairs, amount, etc...)
		var tradeInfo TradeInfo
		tradeInfo, err = svc.getTradeInfo(platform)
		if err != nil {
			return
		}

		// update prices for all pairs and platforms
		err = svc.updateIndexPrice(platform, tradeInfo.PairsToTrade)
		if err != nil {
			return
		}

		// Loop over all pairs to trade
		for _, pair := range tradeInfo.PairsToTrade {
			// get prices history based on interval config
			var prices []model.Price
			priceForm := repository.GetPriceHistoryForm{
				PlatformName: platformName,
				Pair:         pair,
				SinceTime:    time.Now().Add(-time.Duration(svc.algo.PricesNeeded()) * delayBetweenEachRun),
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
			if !ok {
				return fmt.Errorf("cannot get index price for platform %s and pair %s", platformName, pair)
			}
			openForm := tradingStrategy2.ShouldAddOrderForm{
				PriceHistory:       prices,
				IndexPrice:         indexPrice,
				PairToTrade:        pair,
				CurrentBalance:     svc.balanceByPlatform[platformName].ValueByCurrency,
				OpenedOrdersByPair: svc.openedOrdersByPlatformByPair[platformName],
				AllPairsTraded:     initialPairsToTradeByPlatform[platformName],
			}
			var order model.Order
			var shouldOpenOrder bool
			shouldOpenOrder, order, err = svc.algo.ShouldAddOrder(openForm)
			if err != nil {
				return
			}
			// open order if conditions are ok
			if shouldOpenOrder {
				err = svc.addOrder(platform, &order)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (svc *service) updateOpenedOrders(platform tradingPlatform.Api) (err error) {
	if svc.openedOrdersByPlatformByPair[platform.Name()] == nil {
		svc.openedOrdersByPlatformByPair[platform.Name()] = make(map[string][]model.Order)
	}
	orders, err := platform.GetAllOrders()
	if err != nil {
		return
	}
	for _, order := range orders {
		err = svc.repo.StoreOrder(&order)
		if err != nil {
			return
		}
		if order.Status == tradingConst.OpenOrderStatus {
			svc.openedOrdersByPlatformByPair[platform.Name()][order.Pair] = append(
				svc.openedOrdersByPlatformByPair[platform.Name()][order.Pair], order)
		}
	}
	return
}

func (svc *service) getTradeInfo(platform tradingPlatform.Api) (info TradeInfo, err error) {
	// count number of opened orders by pair
	openedOrdersByPlatform, ok := svc.openedOrdersByPlatformByPair[platform.Name()]
	if !ok {
		return info, fmt.Errorf("cannot find opened orders for paltform %s", platform.Name())
	}

	// remove from pairsToTrade all pair that has been reached maxOpenedOrderByPair
	for _, pair := range initialPairsToTradeByPlatform[platform.Name()] {
		if len(openedOrdersByPlatform[pair]) < svc.algo.MaxOpenedOrdersByPair() {
			info.PairsToTrade = append(info.PairsToTrade, pair)
		}
	}

	logger.Infof("Following pairs will be traded : %v", info.PairsToTrade)
	return
}

func (svc *service) updateBalance(platform tradingPlatform.Api) (err error) {
	balance, ok := svc.balanceByPlatform[platform.Name()]
	if !ok {
		return fmt.Errorf("cannot find balance for platform %s", platform.Name())
	}
	err = platform.RefreshBalance(&balance)
	if err != nil {
		return
	}
	svc.balanceByPlatform[platform.Name()] = balance
	err = svc.repo.StoreBalance(&balance)
	if err != nil {
		return
	}

	// init balance if not already done
	if svc.initialBalanceByPlatformByCurrency[platform.Name()] == nil {
		svc.initialBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
		svc.initialBalanceByPlatformByCurrency[platform.Name()] = svc.balanceByPlatform[platform.Name()].ValueByCurrency
	}
	logger.Info(svc.balanceByPlatform[platform.Name()])
	return
}

func (svc *service) updateIndexPrice(platform tradingPlatform.Api, pairs []string) (err error) {
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
			svc.indexPriceByPlatformByPair[platform.Name()] = make(map[string]model.Price)
		}
		svc.indexPriceByPlatformByPair[platform.Name()][price.Pair] = price
		logger.Info(price)
	}
	return
}

func (svc *service) addOrder(platform tradingPlatform.Api, order *model.Order) (err error) {
	// don't open order if balance is less or equal to 0
	if order.Amount <= 0 || order.Volume <= 0 {
		return
	}

	// add order using platform
	err = platform.AddOrder(order)
	if err != nil {
		return
	}
	logger.Info(*order)

	// update balance
	err = svc.updateBalance(platform)
	if err != nil {
		return
	}

	// store order into repository
	err = svc.repo.StoreOrder(order)
	if err != nil {
		return
	}
	return
}
