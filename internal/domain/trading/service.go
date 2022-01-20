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
	getAmountToInvest(platformName string, pair string) (amount float64, err error)
	getAskPrice(platformName string, pair string) (askPrice float64, err error)
}

type service struct {
	quitChannel                        chan bool                     // quit goroutine when program exit
	config                             *Config                       // config that should be applied with trading algorithm
	platformApis                       []tradingPlatform.Api         // communicate with trading platforms
	priceRepo                          repository.Price              // use to store and get previous prices
	positionRepo                       repository.Position           // use to store and get previous positions
	startTime                          time.Time                     // datetime when lago has been started
	initialBalanceByPlatformByCurrency map[string]map[string]float64 // initial balance before opening first position by platform and by currency
	currentBalanceByPlatformByCurrency map[string]map[string]float64 // current balance updated after each run
}

func NewService(config *Config, platformApis []tradingPlatform.Api, priceRepo repository.Price, positionRepo repository.Position) (Service, error) {
	if config == nil {
		defaultConfig, err := GetConfigFromEnvOrDefault()
		if err != nil {
			return nil, err
		}
		config = &defaultConfig
	}
	return &service{
		config:                             config,
		quitChannel:                        make(chan bool),
		platformApis:                       platformApis,
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

	// Update current balance
	for _, platform := range svc.platformApis {
		err = svc.updateWallet(platform.Name())
		if err != nil {
			return
		}
	}

	// Set initial balance before first run
	svc.initialBalanceByPlatformByCurrency = svc.currentBalanceByPlatformByCurrency

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
	for _, platform := range svc.platformApis {
		logger.Infof("-------------------  %s  ------------------------------", platform.Name())

		// update prices for all pairs and platforms
		err = svc.updatePrices(platform.Name())
		if err != nil {
			return
		}

		// Loop over all pairs to trade
		for _, pair := range pairsToTradeByPlatform[platform.Name()] {
			// update wallet for specific platform
			err = svc.updateWallet(platform.Name())
			if err != nil {
				return
			}

			// open new position if conditions are ok
			var shouldOpenPosition bool
			shouldOpenPosition, err = svc.shouldOpenNewPosition(platform.Name(), pair)
			if err != nil {
				return
			}
			if shouldOpenPosition {
				var amount float64
				amount, err = svc.getAmountToInvest(platform.Name(), pair)
				if err != nil {
					return
				}
				if amount >= svc.config.MinimumAmountToOpenPosition {
					var askPrice float64
					askPrice, err = svc.getAskPrice(platform.Name(), pair)
					if err != nil {
						return
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
					ID:       openedPosition.ID,
					Amount:   openedPosition.Amount,
					AskPrice: openedPosition.AskPrice,
				}
				var shouldClosePosition bool
				shouldClosePosition, err = svc.shouldClosePosition(position)
				if err != nil {
					return
				}
				if shouldClosePosition {
					closeForm := tradingPlatform.ClosePositionForm{PositionID: position.ID}
					var closeView tradingPlatform.ClosePositionView
					closeView, err = platform.ClosePosition(closeForm)
					resultInPercent := tradingUtils.GetResultInPercent(closeView.AskPrice, closeView.BidPrice, position.Amount)
					logger.Infof("Closing position %s make profit of %0.2f (%0.2f%%) with ask=%0.2f bid=%0.2f and amount=%0.2f",
						position.ID, closeView.Profit, resultInPercent, closeView.AskPrice, closeView.BidPrice, position.Amount)
				}
			}
		}
	}
	return
}

func (svc *service) shouldOpenNewPosition(platformName string, pair string) (shouldOpen bool, err error) {
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
	estimatedResultInPercent := tradingUtils.GetResultInPercent(position.AskPrice, position.BidPrice, position.Amount)
	if estimatedResultInPercent >= svc.config.MinimumResultInPercentToClosePosition {
		shouldClose = true
	}
	return
}

func (svc *service) updateWallet(platformName string) (err error) {
	for _, platform := range svc.platformApis {
		if platform.Name() != platformName {
			continue // update only wallet for specified platform
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
	}
	logger.Infof("Current wallet is %v", svc.currentBalanceByPlatformByCurrency)
	return
}

func (svc *service) updatePrices(platformName string) (err error) {
	for _, platform := range svc.platformApis {
		if platform.Name() != platformName {
			continue // update prices only for specified platform
		}
		pairsToTrade, ok := pairsToTradeByPlatform[platform.Name()]
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
	}
	return
}

func (svc *service) getAmountToInvest(platformName string, pair string) (amount float64, err error) {
	// Count how many positions has been already opened using this pair for this platform
	openedPositions, err := svc.positionRepo.GetOpenedPositions(platformName)
	if err != nil {
		return
	}
	var nbOpenedPositionsForPair int
	for _, openedPosition := range openedPositions {
		if openedPosition.Pair == pair {
			nbOpenedPositionsForPair++
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
	amount = currentBalance / float64(len(pairsToTradeByPlatform)*svc.config.MaxOpenedPositionsByPair-nbOpenedPositionsForPair)
	return
}

func (svc *service) getAskPrice(platformName string, pair string) (askPrice float64, err error) {
	var lastPrice *repositoryModel.Price
	lastPrice, err = svc.priceRepo.GetCurrentPrice(platformName, pair)
	if err != nil {
		return
	}
	askPrice = lastPrice.AskPrice
	return
}
