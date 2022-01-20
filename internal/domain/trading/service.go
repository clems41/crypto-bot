package trading

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/internal/service/tradingPlatform"
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
	updateWallet(platformName string) (err error)
	updatePrices(platformName string) (err error)
	shouldOpenNewPosition(platformName string, pair string) (shouldOpen bool, err error)
	shouldClosePosition(position Position) (shouldClose bool, err error)
	openPosition(platformName string, pair string) (err error)
	closePosition(position Position) (err error)
	getAmountToInvest(platformName string, pair string) (amount float64, err error)
	getAskPrice(platformName string, pair string) (askPrice float64, err error)
	getBidPrice(platformName string, pair string) (bidPrice float64, err error)
}

type service struct {
	quitChannel                        chan bool                      // quit goroutine when program exit
	config                             *Config                        // config that should be applied with trading algorithm
	platformApis                       map[string]tradingPlatform.Api // communicate with trading platforms
	priceRepo                          repository.Price               // use to store and get previous prices
	positionRepo                       repository.Position            // use to store and get previous positions
	startTime                          time.Time                      // datetime when lago has been started
	initialBalanceByPlatformByCurrency map[string]map[string]float64  // initial balance before opening first position by platform and by currency
	currentBalanceByPlatformByCurrency map[string]map[string]float64  // current balance updated after each run
}

func NewService(config *Config, platformApis []tradingPlatform.Api, priceRepo repository.Price, positionRepo repository.Position) (Service, error) {
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
		priceRepo:                          priceRepo,
		positionRepo:                       positionRepo,
		initialBalanceByPlatformByCurrency: make(map[string]map[string]float64),
		currentBalanceByPlatformByCurrency: make(map[string]map[string]float64),
	}, nil
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()

	// Update current balance and set initial balance
	for _, platform := range svc.platformApis {
		err = svc.updateWallet(platform.Name())
		if err != nil {
			return
		}
		for currency, balance := range svc.currentBalanceByPlatformByCurrency[platform.Name()] {
			if svc.initialBalanceByPlatformByCurrency[platform.Name()] == nil {
				svc.initialBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
			}
			svc.initialBalanceByPlatformByCurrency[platform.Name()][currency] = balance
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

	// Close all positions for all platforms
	for _, platform := range svc.platformApis {
		logger.Infof("-------------------  %s  ------------------------------", platform.Name())
		// Close opened positions
		var openedPositions tradingPlatform.GetOpenedPositionsView
		openedPositions, err = platform.GetOpenedPositions()
		if err != nil {
			return
		}
		for _, position := range openedPositions.Positions {
			form := tradingPlatform.ClosePositionForm{PositionID: position.ID}
			_, err = platform.ClosePosition(form)
			if err != nil {
				return
			}
		}

		// Get all positions
		var positions tradingPlatform.GetAllPositionsView
		positions, err = platform.GetAllPositions()
		if err != nil {
			return
		}
		for _, position := range positions.Positions {
			logger.Infof("Position %s pair=%s ask=%0.2f bid=%0.2f make result of %0.4f with amount %0.2f (%0.3f%%)",
				position.ID, position.Pair, position.AskPrice, position.BidPrice, position.Result, position.Amount,
				tradingUtils.GetResultInPercent(position.AskPrice, position.BidPrice, position.Amount))
		}

		// Update current balance
		err = svc.updateWallet(platform.Name())
		if err != nil {
			return
		}

		// Calculate estimated profit
		endTime := time.Now()
		tradingDuration := endTime.Sub(svc.startTime)
		var initialBalance, finalBalance float64
		for currency, balance := range svc.currentBalanceByPlatformByCurrency[platform.Name()] {
			_, ok := svc.initialBalanceByPlatformByCurrency[platform.Name()]
			if !ok {
				return errCurrencyNotInBalance
			}
			initialBalanceCurrency, ok := svc.initialBalanceByPlatformByCurrency[platform.Name()][currency]
			if !ok {
				return errCurrencyNotInBalance
			}
			initialBalance += initialBalanceCurrency
			finalBalance += balance
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
	logger.Infof("Current wallet is %v", svc.currentBalanceByPlatformByCurrency)
	for platformName, platform := range svc.platformApis {
		logger.Infof("--  %s  --", platformName)

		// update prices for all pairs and platforms
		err = svc.updatePrices(platformName)
		if err != nil {
			return
		}

		// Loop over all pairs to trade
		for _, pair := range pairsToTradeByPlatform[platformName] {
			// open new position if conditions are ok
			var shouldOpenPosition bool
			shouldOpenPosition, err = svc.shouldOpenNewPosition(platformName, pair)
			if err != nil {
				return
			}
			if shouldOpenPosition {
				err = svc.openPosition(platformName, pair)
				if err != nil {
					return
				}
			}
		}

		// close positions if condition are ok
		var openedPositionsView tradingPlatform.GetOpenedPositionsView
		openedPositionsView, err = platform.GetOpenedPositions()
		if err != nil {
			return
		}
		for _, openedPosition := range openedPositionsView.Positions {
			position := Position{
				ID:           openedPosition.ID,
				PlatformName: platformName,
				Amount:       openedPosition.Amount,
				AskPrice:     openedPosition.AskPrice,
				Pair:         openedPosition.Pair,
			}
			var shouldClosePosition bool
			shouldClosePosition, err = svc.shouldClosePosition(position)
			if err != nil {
				return
			}
			if shouldClosePosition {
				err = svc.closePosition(position)
			}
		}
	}
	return
}

func (svc *service) shouldOpenNewPosition(platformName string, pair string) (shouldOpen bool, err error) {
	// If too many positions has been opened for this platform/pair, don't open new one
	openedPositions, err := svc.positionRepo.GetOpenedPositions(platformName, pair)
	if err != nil {
		return
	}
	if len(openedPositions) >= svc.config.MaxOpenedPositionsByPair {
		return false, nil
	}

	// Getting last prices for current pair
	lastPrices, err := svc.priceRepo.GetLast(platformName, svc.config.NumberOfPreviousPricesToCompare, pair)
	if err != nil {
		return
	}
	if len(lastPrices) < svc.config.NumberOfPreviousPricesToCompare {
		return false, errPriceNotFound
	}

	// Find minimum value from last prices
	minimumPrice := lastPrices[0].AskPrice
	for _, price := range lastPrices {
		if price.AskPrice < minimumPrice {
			minimumPrice = price.AskPrice
		}
	}

	// Position should be opened if last price is the lowest of all last 10 prices
	lastPrice := lastPrices[len(lastPrices)-1]
	if minimumPrice == lastPrice.AskPrice {
		shouldOpen = true
	}
	return
}

func (svc *service) shouldClosePosition(position Position) (shouldClose bool, err error) {
	bidPrice, err := svc.getBidPrice(position.PlatformName, position.Pair)
	if err != nil {
		return
	}
	estimatedResultInPercent := tradingUtils.GetResultInPercent(position.AskPrice, bidPrice, position.Amount)
	if estimatedResultInPercent >= svc.config.MinimumResultInPercentToClosePosition {
		shouldClose = true
	}
	return
}

func (svc *service) updateWallet(platformName string) (err error) {
	platform, ok := svc.platformApis[platformName]
	if !ok {
		return errPlatformNotFound
	}
	var walletView tradingPlatform.WalletView
	walletView, err = platform.GetWalletBalance()
	if err != nil {
		return
	}
	for currency, balance := range walletView.BalanceByCurrency {
		if svc.currentBalanceByPlatformByCurrency[platform.Name()] == nil {
			svc.currentBalanceByPlatformByCurrency[platform.Name()] = make(map[string]float64)
		}
		svc.currentBalanceByPlatformByCurrency[platform.Name()][currency] = balance
	}
	logger.Infof("Current wallet is %v", svc.currentBalanceByPlatformByCurrency)
	return
}

func (svc *service) updatePrices(platformName string) (err error) {
	platform, ok := svc.platformApis[platformName]
	if !ok {
		return errPlatformNotFound
	}
	pairsToTrade, ok := pairsToTradeByPlatform[platformName]
	if !ok {
		return errPairNotFound
	}
	priceForm := tradingPlatform.GetPriceForm{
		Pairs: pairsToTrade,
	}
	var priceView tradingPlatform.GetPriceView
	priceView, err = platform.GetPrice(priceForm)
	if err != nil {
		return
	}
	for pair, price := range priceView.PriceByPair {
		// Store new prices (use after to know if new open position should be opened
		priceModel := repositoryModel.Price{
			Date:         price.Date,
			Pair:         pair,
			AskPrice:     price.AskPrice,
			BidPrice:     price.BidPrice,
			PlatformName: platform.Name(),
		}
		err = svc.priceRepo.Store(&priceModel)
		if err != nil {
			return
		}
		logger.Infof("%s price is ask=%0.2f and bid=%0.2f", pair, price.AskPrice, price.BidPrice)
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
	currency := tradingUtils.CurrencyNeededToTradePair(pair)
	for _, pairTraded := range pairsToTradeByPlatform[platformName] {
		currencyNeededForPair := tradingUtils.CurrencyNeededToTradePair(pairTraded)
		if currencyNeededForPair == currency {
			nbPairUsingSameCurrency++
		}
	}

	// Get current balance for currency
	currentBalance, ok := svc.currentBalanceByPlatformByCurrency[platformName][currency]
	if !ok {
		return amount, errCurrencyNotInBalance
	}

	// Calculate amount to invest
	amount = currentBalance / float64(len(pairsToTradeByPlatform[platformName])*svc.config.MaxOpenedPositionsByPair-nbOpenedPositionsForPlatform)
	return
}

func (svc *service) getAskPrice(platformName string, pair string) (askPrice float64, err error) {
	var lastPrice *repositoryModel.Price
	lastPrice, err = svc.priceRepo.GetCurrentPrice(platformName, pair)
	if err != nil {
		return
	}
	if lastPrice == nil {
		return askPrice, errPriceNotFound
	}
	askPrice = lastPrice.AskPrice
	return
}

func (svc *service) getBidPrice(platformName string, pair string) (bidPrice float64, err error) {
	var lastPrice *repositoryModel.Price
	lastPrice, err = svc.priceRepo.GetCurrentPrice(platformName, pair)
	if err != nil {
		return
	}
	if lastPrice == nil {
		return bidPrice, errPriceNotFound
	}
	bidPrice = lastPrice.BidPrice
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
		position := repositoryModel.Position{
			ID:           openView.PositionID,
			PlatformName: platformName,
			AskDate:      time.Now(),
			//BidDate:      time.Time{},
			Pair:     openForm.Pair,
			Amount:   openView.Amount,
			AskPrice: openView.AskPrice,
			//BidPrice:     0,
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
}
