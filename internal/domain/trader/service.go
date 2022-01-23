package trader

import (
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/internal/service/tradingStrategy"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/tradingUtils"
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
	addOrder(platform tradingPlatform.Api, pair string) (err error)
}

type service struct {
	quitChannel                        chan bool                          // quit goroutine when program exit
	config                             *Config                            // config that should be applied with trading algorithm
	platformApis                       map[string]tradingPlatform.Api     // communicate with trading platforms
	repo                               repository.Repository              // use to store data
	algo                               tradingStrategy.Algo               // use to know if order should be open based on prices
	startTime                          time.Time                          // datetime when lago has been started
	initialBalanceByPlatformByCurrency map[string]map[string]float64      // initial balance before opening first position by platform and by currency
	currentBalanceByPlatformByCurrency map[string]map[string]float64      // current balance
	indexPriceByPlatformByPair         map[string]map[string]*model.Price // store last index got from platform
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
		currentBalanceByPlatformByCurrency: make(map[string]map[string]float64),
		indexPriceByPlatformByPair:         make(map[string]map[string]*model.Price),
	}, nil
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()

	// Update current balance
	for _, platform := range svc.platformApis {
		err = svc.updateBalance(platform)
		if err != nil {
			return
		}
	}

	// Set initial balance, useful to calculate ending profit, result, etc...
	for _, platform := range svc.platformApis {
		var balance *model.Balance
		balance, err = platform.GetBalance()
		if err != nil {
			return
		}
		for currency, value := range balance.ValueByCurrency {
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
		var orders []*model.Order
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
			finalBalanceCurrency, ok := svc.currentBalanceByPlatformByCurrency[platformName][currency]
			if !ok {
				return errCurrencyNotInBalance(currency)
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
			// open new position if conditions are ok
			var shouldAddOrder bool
			shouldAddOrder, err = svc.algo.ShouldAddOrder()
			if err != nil {
				return
			}
			if shouldAddOrder {
				err = svc.addOrder(platform, pair)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (svc *service) updateBalance(platform tradingPlatform.Api) (err error) {
	balance, err := platform.GetBalance()
	if err != nil {
		return
	}
	err = svc.repo.StoreBalance(balance)
	if err != nil {
		return
	}
	if svc.currentBalanceByPlatformByCurrency[platform.Name()] == nil {
		svc.currentBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
	}
	svc.currentBalanceByPlatformByCurrency[platform.Name()] = balance.ValueByCurrency
	logger.Infof("Platform %s --> Balance %v at %s", platform.Name(), balance.ValueByCurrency,
		balance.UpdatedAt.Format(time.RFC3339))
	return
}

func (svc *service) updateIndexPrice(platform tradingPlatform.Api) (err error) {
	pairs, ok := pairsToTradeByPlatform[platform.Name()]
	if !ok {
		return errPlatformNotExist(platform.Name())
	}
	prices, err := platform.GetIndexPrices(pairs...)
	if err != nil {
		return
	}
	for pair, price := range prices {
		err = svc.repo.StorePrice(price)
		if err != nil {
			return
		}
		if svc.indexPriceByPlatformByPair[platform.Name()] == nil {
			svc.indexPriceByPlatformByPair[platform.Name()] = make(map[string]*model.Price)
		}
		svc.indexPriceByPlatformByPair[platform.Name()][pair] = price
		logger.Infof("Platform %s --> %s ask=%0.2f bid=%0.2f at %s", platform.Name(), pair, price.Ask,
			price.Bid, price.Date.Format(time.RFC3339))
	}
	return
}

func (svc *service) addOrder(platform tradingPlatform.Api, pair string) (err error) {
	// creating order
	order := model.Order{
		Side:                "",
		Volume:              0,
		Type:                "",
		Price:               0,
		Amount:              0,
		Leverage:            0,
		CloseConditionType:  "",
		CloseConditionPrice: 0,
		Fees:                0,
		Status:              "",
	}

	// add order using platform
	err = platform.AddOrder(&order)
	if err != nil {
		return
	}

	// store order into repository
	err = svc.repo.StoreOrder(&order)
	if err != nil {
		return
	}
	logger.Infof("New order has been added %v", order)
	return
}

/*func (svc *service) openPosition(platformName string, pair string) (err error) {
	// Get amount to invest in new position
	var amount float64
	amount, err = svc.getAmountToInvest(platformName, pair)
	if err != nil {
		return
	}

	// open only if amount is more than minimum defined in config
	if amount >= svc.config.MinimumAmountToOpenPosition {
		// get ask price
		var askPrice float64
		askPrice, err = svc.getAskPrice(platformName, pair)
		if err != nil {
			return
		}
		platform, ok := svc.platformApis[platformName]
		if !ok {
			return errPlatformNotFound
		}
		openForm := tradingPlatform.OpenPositionForm{
			Pair:     pair,
			Amount:   amount,
			AskPrice: askPrice,
		}
		var openView tradingPlatform.OpenPositionView
		openView, err = platform.OpenPosition(openForm)
		if err != nil {
			return
		}
		logger.Infof("Opening new position %s for pair %s with ask=%0.2f and amount=%0.2f",
			openView.PositionID, openForm.Pair, openView.AskPrice, openView.Amount)

		// Store position in repo for KPI
		expectedBidPrice := (1 + svc.config.MinimumResultInPercentToClosePosition/100) * openView.AskPrice
		position := model.Position{
			ID:           openView.PositionID,
			PlatformName: platformName,
			AskDate:      time.Now(),
			//BidDate:      time.Time{},
			Pair:     openForm.Pair,
			Amount:   openView.Amount,
			AskPrice: openView.AskPrice,
			//Bid:     0,
			//Result:       0,
			ExpectedBidPrice: expectedBidPrice,
			Closed:           false,
		}
		err = svc.positionRepo.Store(&position)
		if err != nil {
			return
		}

		// When position has been opened, wallet should be updated
		err = svc.updateWallet(platformName)
		if err != nil {
			return
		}
	}
	return
}


func (svc *service) getAmountToInvest(platformName string, pair string) (amount float64, err error) {
	// Count how many positions has been already opened using this pair for this platform
	openedPositions, err := svc.positionRepo.GetOpenedPositions(platformName)
	if err != nil {
		return
	}
	var nbOpenedPositionsForPair, nbOpenedPositionsForPlatform int
	for _, openedPosition := range openedPositions {
		if openedPosition.PlatformName == platformName {
			nbOpenedPositionsForPlatform++
			if openedPosition.Pair == pair {
				nbOpenedPositionsForPair++
			}
		}
	}
	if nbOpenedPositionsForPair == svc.config.MaxOpenedPositionsByPair {
		return
	}

	// Get number of pair that can be trade with same currency
	var nbPairUsingSameCurrency int
	currency := CurrencyNeededToTradePair(pair)
	for _, pairTraded := range pairsToTradeByPlatform[platformName] {
		currencyNeededForPair := CurrencyNeededToTradePair(pairTraded)
		if currencyNeededForPair == currency {
			nbPairUsingSameCurrency++
		}
	}

	// Get current balance for currency
	currentBalance, err := svc.balanceRepo.GetCurrentBalance(platformName, currency)
	if err != nil {
		return
	}

	// Calculate amount to invest
	amount = currentBalance.Value / float64(len(pairsToTradeByPlatform[platformName])*svc.config.MaxOpenedPositionsByPair-nbOpenedPositionsForPlatform)
	return
}

func (svc *service) getAskPrice(platformName string, pair string) (askPrice float64, err error) {
	var lastPrice *model.Price
	lastPrice, err = svc.priceRepo.GetCurrentPrice(platformName, pair)
	if err != nil {
		return
	}
	if lastPrice == nil {
		return askPrice, errPriceNotFound
	}
	askPrice = lastPrice.Ask
	return
}

func (svc *service) getBidPrice(platformName string, pair string) (bidPrice float64, err error) {
	var lastPrice *model.Price
	lastPrice, err = svc.priceRepo.GetCurrentPrice(platformName, pair)
	if err != nil {
		return
	}
	if lastPrice == nil {
		return bidPrice, errPriceNotFound
	}
	bidPrice = lastPrice.Bid
	return
}

func (svc *service) openPosition(platformName string, pair string) (err error) {
	// Get amount to invest in new position
	var amount float64
	amount, err = svc.getAmountToInvest(platformName, pair)
	if err != nil {
		return
	}

	// open only if amount is more than minimum defined in config
	if amount >= svc.config.MinimumAmountToOpenPosition {
		// get ask price
		var askPrice float64
		askPrice, err = svc.getAskPrice(platformName, pair)
		if err != nil {
			return
		}
		platform, ok := svc.platformApis[platformName]
		if !ok {
			return errPlatformNotFound
		}
		openForm := tradingPlatform.OpenPositionForm{
			Pair:     pair,
			Amount:   amount,
			AskPrice: askPrice,
		}
		var openView tradingPlatform.OpenPositionView
		openView, err = platform.OpenPosition(openForm)
		if err != nil {
			return
		}
		logger.Infof("Opening new position %s for pair %s with ask=%0.2f and amount=%0.2f",
			openView.PositionID, openForm.Pair, openView.AskPrice, openView.Amount)

		// Store position in repo for KPI
		expectedBidPrice := (1 + svc.config.MinimumResultInPercentToClosePosition/100) * openView.AskPrice
		position := model.Position{
			ID:           openView.PositionID,
			PlatformName: platformName,
			AskDate:      time.Now(),
			//BidDate:      time.Time{},
			Pair:     openForm.Pair,
			Amount:   openView.Amount,
			AskPrice: openView.AskPrice,
			//Bid:     0,
			//Result:       0,
			ExpectedBidPrice: expectedBidPrice,
			Closed:           false,
		}
		err = svc.positionRepo.Store(&position)
		if err != nil {
			return
		}

		// When position has been opened, wallet should be updated
		err = svc.updateWallet(platformName)
		if err != nil {
			return
		}
	}
	return
}

func (svc *service) closePosition(position Position) (err error) {
	closeForm := tradingPlatform.ClosePositionForm{
		PositionID: position.ID,
	}
	var closeView tradingPlatform.ClosePositionView
	platform, ok := svc.platformApis[position.PlatformName]
	if !ok {
		return errPlatformNotFound
	}
	closeView, err = platform.ClosePosition(closeForm)
	resultInPercent := tradingUtils.GetResultInPercent(closeView.AskPrice, closeView.BidPrice, closeView.Amount)
	result := tradingUtils.GetResult(closeView.AskPrice, closeView.BidPrice, closeView.Amount)
	profit := tradingUtils.GetProfit(closeView.AskPrice, closeView.BidPrice, closeView.Amount)
	logger.Infof("Closing position %s make profit of %0.2f (%0.2f%%) with ask=%0.2f bid=%0.2f and amount=%0.2f",
		position.ID, profit, resultInPercent, closeView.AskPrice, closeView.BidPrice, position.Amount)

	// When position has been closed, wallet should be updated
	err = svc.updateWallet(position.PlatformName)
	if err != nil {
		return
	}

	// Update position from repo, useful for KPI
	positionToUpdate, err := svc.positionRepo.Get(position.ID)
	if err != nil {
		return
	}
	positionToUpdate.BidPrice = closeView.BidPrice
	positionToUpdate.BidDate = time.Now()
	positionToUpdate.Closed = true
	positionToUpdate.Result = result
	positionToUpdate.ResultInPercent = resultInPercent
	positionToUpdate.Profit = profit
	err = svc.positionRepo.Store(positionToUpdate)
	if err != nil {
		return
	}

	return
}*/
