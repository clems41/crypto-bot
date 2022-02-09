package trader

import (
	"crypto-bot/external/service/mailService"
	"crypto-bot/external/service/tradingPlatform"
	"crypto-bot/internal/constant/timeConst"
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/domain/tradingStrategy"
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
	Stop()

	/* Private methods */
	applyTradingAlgorithm() (err error)
	getTradeInfo(platform tradingPlatform.Api) (info TradeInfo, err error)
	updateOpenedOrders(platform tradingPlatform.Api) (err error)
	updateBalance(platform tradingPlatform.Api) (err error)
	updateIndexPrice(platform tradingPlatform.Api, pairs []tradingConst.Pair) (err error)
	addOrder(platform tradingPlatform.Api, order *model.Order) (err error)
}

type service struct {
	quitChannel                        chan bool                                      // quit goroutine when program exit
	platformApis                       map[string]tradingPlatform.Api                 // communicate with trading platforms
	repo                               repository.Repository                          // use to store data
	algo                               tradingStrategy.Algo                           // use to know if order should be open based on prices
	startTime                          time.Time                                      // datetime when lago has been started
	mailService                        mailService.Service                            // use to send order by mail
	initialBalanceByPlatformByCurrency map[string]map[tradingConst.Currency]float64   // initial balance before opening first order by platform and by currency
	balanceByPlatform                  map[string]model.Balance                       // current balance
	indexPriceByPlatformByPair         map[string]map[tradingConst.Pair]model.Price   // store last index got from platform
	openedOrdersByPlatformByPair       map[string]map[tradingConst.Pair][]model.Order // store opened orders calculating amount to invest by pair
	previousOrdersByPlatformById       map[string]map[string]model.Order              // store orders
}

func NewService(platformApis []tradingPlatform.Api, repo repository.Repository, algo tradingStrategy.Algo,
	mailService mailService.Service) (Service, error) {
	platformApisMap := make(map[string]tradingPlatform.Api)
	for _, platformApi := range platformApis {
		platformApisMap[platformApi.Name()] = platformApi
	}
	return &service{
		quitChannel:                        make(chan bool),
		platformApis:                       platformApisMap,
		repo:                               repo,
		algo:                               algo,
		mailService:                        mailService,
		initialBalanceByPlatformByCurrency: make(map[string]map[tradingConst.Currency]float64),
		balanceByPlatform:                  make(map[string]model.Balance),
		indexPriceByPlatformByPair:         make(map[string]map[tradingConst.Pair]model.Price),
		openedOrdersByPlatformByPair:       make(map[string]map[tradingConst.Pair][]model.Order),
		previousOrdersByPlatformById:       make(map[string]map[string]model.Order),
	}, nil
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()

	// init all maps
	for platformName, platform := range svc.platformApis {
		svc.balanceByPlatform[platformName] = model.Balance{
			PlatformName:    platformName,
			ValueByCurrency: make(map[tradingConst.Currency]float64),
		}
		svc.previousOrdersByPlatformById[platformName] = make(map[string]model.Order)
		svc.openedOrdersByPlatformByPair[platformName] = make(map[tradingConst.Pair][]model.Order)
		svc.initialBalanceByPlatformByCurrency[platformName] = make(map[tradingConst.Currency]float64)
		svc.indexPriceByPlatformByPair[platformName] = make(map[tradingConst.Pair]model.Price)

		// Set initial balance, useful to calculate ending profit, result, etc...
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}
		for currency, value := range svc.balanceByPlatform[platformName].ValueByCurrency {
			svc.initialBalanceByPlatformByCurrency[platformName][currency] = value
		}
	}

	// Run algorithm each X ms
	for range time.Tick(svc.algo.DelayBetweenEachRun()) { // Loop
		select {
		case <-svc.quitChannel:
			logger.Debugf("Stop message has been received")
			return
		default:
			err = svc.applyTradingAlgorithm()
			if err != nil {
				sendRequest := mailService.SendRequest{
					Subject: "Error occurs with Crypto-bot",
					Body:    err.Error(),
				}
				errMail := svc.mailService.Send(sendRequest)
				if errMail != nil {
					logger.Error(errMail)
				}
				return
			}
		}
	}
	return
}

func (svc *service) Stop() {
	svc.quitChannel <- true

	// Cancel all open orders
	for platformName, platform := range svc.platformApis {
		logger.Infof("--  %s  --", platformName)
		// Get all orders
		orders, err := platform.GetAllOrders(svc.startTime)
		if err != nil {
			return
		}
		for _, order := range orders {
			logger.Info(order)
		}

		// Get current balance
		err = svc.updateBalance(platform)
		if err != nil {
			logger.Error(err)
			return
		}

		svc.calculateEndProfit(platform)
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

		// update prices for all pairs and platforms
		pairs, ok := svc.algo.PairsToTradeByPlatform()[platformName]
		if !ok {
			return fmt.Errorf("cannot find pairs to trade for platform %s", platformName)
		}
		err = svc.updateIndexPrice(platform, pairs)
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

		// Loop over all pairs to trade
		for _, pair := range tradeInfo.PairsToTrade {
			// get prices history based on interval config
			var prices []model.Price
			priceForm := repository.GetPriceHistoryForm{
				PlatformName: platformName,
				Pair:         pair,
				SinceTime:    time.Now().Add(-time.Duration(svc.algo.PricesNeeded()) * svc.algo.DelayBetweenEachRun()),
			}
			prices, err = svc.repo.GetPriceHistory(priceForm)
			if err != nil {
				return
			}

			// If price history doesn't return enough prices, svc.algo.ShouldAddOrder will return an error.
			// So we skip this run, and try again the next one.
			// It can take some time to get enough prices at start, depending on delayBetweenEachRun.
			if len(prices) < svc.algo.PricesNeeded() {
				logger.Debugf("Run for pair %s will be skip, doesn't get enough prices from history (need=%d got=%d)",
					pair, svc.algo.PricesNeeded(), len(prices))
				continue
			}

			// find order history
			orderForm := repository.GetOrderHistoryForm{
				PlatformName: platformName,
				Status:       tradingConst.Open,
			}
			var openOrders []model.Order
			openOrders, err = svc.repo.GetOrderHistory(orderForm)
			if err != nil {
				return
			}

			// fill open form to know if new order should be open
			indexPrice, ok := svc.indexPriceByPlatformByPair[platformName][pair]
			if !ok {
				return fmt.Errorf("cannot get index price for platform %s and pair %s", platformName, pair)
			}
			openForm := tradingStrategy.ShouldAddOrderForm{
				PriceHistory:       prices,
				IndexPrice:         indexPrice,
				PairToTrade:        pair,
				CurrentBalance:     svc.balanceByPlatform[platformName].ValueByCurrency,
				OpenOrders:         openOrders,
				PlatformName:       platformName,
				TakerFeesInPercent: platform.TakerFeesInPercent(),
				MakerFeesInPercent: platform.MakerFeesInPercent(),
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
				break // open order only once at run, avoid updating tradeInfo after each new order
			}
		}
	}
	return
}

func (svc *service) updateOpenedOrders(platform tradingPlatform.Api) (err error) {
	// reset opened orders
	svc.openedOrdersByPlatformByPair[platform.Name()] = make(map[tradingConst.Pair][]model.Order)

	orders, err := platform.GetAllOrders(svc.startTime)
	if err != nil {
		return
	}
	for _, order := range orders {
		err = svc.repo.StoreOrder(&order)
		if err != nil {
			return
		}

		// fill open orders map
		if order.Status == tradingConst.Open {
			svc.openedOrdersByPlatformByPair[platform.Name()][order.Pair] = append(
				svc.openedOrdersByPlatformByPair[platform.Name()][order.Pair], order)
		}

		// send order by mail if it has just been closed
		previousOrder, ok := svc.previousOrdersByPlatformById[platform.Name()][order.ID]
		if !ok || previousOrder.Status != order.Status {
			// if new order or status has been updated, send email with order info
			sendRequest := mailService.SendRequest{
				Subject: fmt.Sprintf("New order from %s", platform.Name()),
				Body:    order.String(),
			}
			err = svc.mailService.Send(sendRequest)
			if err != nil {
				return
			}
		}
		svc.previousOrdersByPlatformById[platform.Name()][order.ID] = order
	}
	return
}

func (svc *service) getTradeInfo(platform tradingPlatform.Api) (info TradeInfo, err error) {
	info.NbOpenOrderByPair = make(map[tradingConst.Pair]int)
	// count number of opened orders by pair
	openedOrdersByPlatform, ok := svc.openedOrdersByPlatformByPair[platform.Name()]
	if !ok {
		return info, fmt.Errorf("cannot find opened orders for paltform %s", platform.Name())
	}

	// remove from pairsToTrade all pair that has been reached maxOpenedOrderByPair
	initialPairs, ok := svc.algo.PairsToTradeByPlatform()[platform.Name()]
	if !ok {
		return info, fmt.Errorf("cannot find pairs to trade for platform %s", platform.Name())
	}
	for _, pair := range initialPairs {
		if len(openedOrdersByPlatform[pair]) < svc.algo.MaxOpenedOrdersByPair() {
			info.PairsToTrade = append(info.PairsToTrade, pair)
		}
		for _, order := range openedOrdersByPlatform[pair] {
			if order.Status == tradingConst.Open {
				info.NbOpenOrderByPair[pair]++
			}
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
	logger.Info(svc.balanceByPlatform[platform.Name()])
	return
}

func (svc *service) updateIndexPrice(platform tradingPlatform.Api, pairs []tradingConst.Pair) (err error) {
	prices, err := platform.GetIndexPrices(pairs...)
	if err != nil {
		return
	}
	for _, price := range prices {
		err = svc.repo.StorePrice(&price)
		if err != nil {
			return
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

	// fill missing order fields
	order.PlatformName = platform.Name()

	// add order using platform
	err = platform.AddOrder(*order)
	if err != nil {
		return
	}
	logger.Info(order)

	// update balance
	err = svc.updateBalance(platform)
	if err != nil {
		return
	}
	return
}

func (svc *service) calculateEndProfit(platform tradingPlatform.Api) {
	// Calculate estimated profit
	endTime := time.Now()
	tradingDuration := endTime.Sub(svc.startTime)
	for currency, initialBalanceCurrency := range svc.initialBalanceByPlatformByCurrency[platform.Name()] {
		finalBalanceCurrency, ok := svc.balanceByPlatform[platform.Name()].ValueByCurrency[currency]
		if !ok {
			logger.Errorf("cannot find balance for currency %s", currency)
			return
		}
		oneDayProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 24*time.Hour)
		oneMonthProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 30*24*time.Hour)
		oneYearProfit := tradingUtils.EstimateProfit(initialBalanceCurrency, finalBalanceCurrency, tradingDuration, 365*24*time.Hour)
		logger.Infof("With initial balance for currency %s of %f, profit for one day would be %f, for one month %f and for one year %f",
			currency, initialBalanceCurrency, oneDayProfit, oneMonthProfit, oneYearProfit)
	}
}
