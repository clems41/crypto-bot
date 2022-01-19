package trading

import (
	"crypto-bot/internal/constant/tradingConst"
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
}

type service struct {
	quitChannel                       chan bool                     // quit goroutine when program exit
	platformApis                      []tradingPlatform.Api         // communicate with trading platforms
	priceRepository                   repository.Price              // use to store and get previous prices
	positionRepo                      repository.Position           // use to store and get previous positions
	startTime                         time.Time                     // datetime when lago has been started
	initialWalletByPlatformByCurrency map[string]map[string]float64 // initial wallet before opening first position by platform and by currency
	openedPositionIdsByTimestamp      map[int64]string              // use to know if position has been opened on specific timestamp
}

func NewService(platformApis []tradingPlatform.Api,
	priceRepo repository.Price,
	positionRepo repository.Position) Service {
	return &service{
		quitChannel:                       make(chan bool),
		platformApis:                      platformApis,
		priceRepository:                   priceRepo,
		positionRepo:                      positionRepo,
		openedPositionIdsByTimestamp:      make(map[int64]string),
		initialWalletByPlatformByCurrency: make(map[string]map[string]float64),
	}
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
	svc.startTime = time.Now()
	for range time.Tick(delayBeforeNewAlgoApplicationInMilliseconds * time.Millisecond) { // Loop
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
		endTime := time.Now()
		tradingDuration := endTime.Sub(svc.startTime)
		var walletView tradingPlatform.WalletView
		walletView, err = platform.GetWalletBalance()
		if err != nil {
			return
		}
		logger.Infof("Current balance is now of %v after %v of trading", walletView.BalanceByCurrency, tradingDuration)

		// Calculate estimated profit
		var initialBalance, finalBalance float64
		for currency, balance := range walletView.BalanceByCurrency {
			_, ok := svc.initialWalletByPlatformByCurrency[platform.Name()]
			if !ok {
				return errCurrencyNotInBalance
			}
			initialBalanceCurrency, ok := svc.initialWalletByPlatformByCurrency[platform.Name()][currency]
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
		// Update current balance
		var walletView tradingPlatform.WalletView
		walletView, err = platform.GetWalletBalance()
		if err != nil {
			return
		}
		logger.Infof("Current balance is %v", walletView.BalanceByCurrency)
		if _, ok := svc.initialWalletByPlatformByCurrency[platform.Name()]; !ok { // if not init, it means it's the first run, so update it
			svc.initialWalletByPlatformByCurrency[platform.Name()] = make(map[string]float64)
			for currency, balance := range walletView.BalanceByCurrency {
				svc.initialWalletByPlatformByCurrency[platform.Name()][currency] = balance
			}
		}

		// Getting actual price of pair
		priceForm := tradingPlatform.GetPriceForm{
			Pairs: pairsToTradeByPlatform[platform.Name()],
		}
		var priceView tradingPlatform.GetPriceView
		priceView, err = platform.GetPrice(priceForm)
		if err != nil {
			return
		}

		// Store new prices
		for pair, price := range priceView.PriceByPair {
			priceModel := repositoryModel.Price{
				Date:     price.Date,
				Pair:     pair,
				AskPrice: price.AskPrice,
				BidPrice: price.BidPrice,
			}
			err = svc.priceRepository.Store(&priceModel)
			if err != nil {
				return
			}
			logger.Infof("%s price is ask=%0.2f and bid=%0.2f", pair, price.AskPrice, price.BidPrice)
		}

		// Loop over all pairs to trade
		for _, pair := range pairsToTradeByPlatform[platform.Name()] {
			// Get pair price
			price, ok := priceView.PriceByPair[pair]
			if !ok {
				return errPairNotFound
			}

			// Get currency needed to trade this pair
			currency := tradingConst.CurrencyNeededToTradePair(pair)
			balanceForCurrency, ok := walletView.BalanceByCurrency[currency]
			if !ok {
				return errCurrencyNotInBalance
			}

			// Get opened positions
			var openedPositions tradingPlatform.GetOpenedPositionsView
			openedPositions, err = platform.GetOpenedPositions()
			if err != nil {
				return
			}

			// Open and close positions depending on new prices and previous positions
			var nbOpenedPositionsForCurrency int
			for _, position := range openedPositions.Positions {
				if tradingConst.CanTradePairUsingCurrency(position.Pair, currency) {
					nbOpenedPositionsForCurrency++
				}
			}
			if nbOpenedPositionsForCurrency >= maxOpenedPositionsByCurrency {
				logger.Infof("No new position for pair %s will be opened because max %d has been reached",
					pair, maxOpenedPositionsByCurrency)
			} else if balanceForCurrency > 0 {
				err = svc.openNewPositions(balanceForCurrency, pair, nbOpenedPositionsForCurrency, platform)
				if err != nil {
					return
				}
			}

			err = svc.closePositions(pair, price.BidPrice, openedPositions.Positions, platform)
			if err != nil {
				return
			}
		}
	}
	return
}

func (svc *service) openNewPositions(actualBalance float64, pair string, nbOpenedPosition int, platform tradingPlatform.Api) (err error) {
	// Get last prices
	var lastPrices []*repositoryModel.Price
	lastPrices, err = svc.priceRepository.GetLast(nbOfIncreasingValueToOpenPosition+1, pair)
	if err != nil {
		return
	}

	shouldOpenPosition, err := svc.shouldOpenPosition(actualBalance, lastPrices, nbOpenedPosition)
	if err != nil {
		return
	}

	if shouldOpenPosition {
		lastPrice := lastPrices[len(lastPrices)-1]
		form := tradingPlatform.OpenPositionForm{
			AskPrice: lastPrice.AskPrice,
			Pair:     pair,
			Amount:   svc.getAmountToInvest(actualBalance, nbOpenedPosition),
		}
		var view tradingPlatform.OpenPositionView
		view, err = platform.OpenPosition(form)
		if err != nil {
			return
		}
		logger.Infof("Opening new position %s with amount of %f for pair %s", view.PositionID, form.Amount, pair)
		svc.openedPositionIdsByTimestamp[lastPrice.Date.Unix()] = view.PositionID
		position := repositoryModel.Position{
			ID:       view.PositionID,
			AskDate:  time.Now(),
			Pair:     form.Pair,
			Amount:   form.Amount,
			AskPrice: form.AskPrice,
			Closed:   false,
		}
		err = svc.positionRepo.Store(&position)
		if err != nil {
			return
		}
	}
	return
}

func (svc *service) closePositions(pair string, actualBidPrice float64, openedPositions []tradingPlatform.PositionView, platform tradingPlatform.Api) (err error) {
	// close position if actual price is more than 0.1% of price position
	for _, position := range openedPositions {
		if position.Pair != pair {
			continue // skip if pair position is not the same to avoid miscalculation
		}
		resultInPercent := tradingUtils.GetResultInPercent(position.AskPrice, actualBidPrice, position.Amount)
		if resultInPercent > resultToClosePositionInPercent {
			form := tradingPlatform.ClosePositionForm{PositionID: position.ID}
			var view tradingPlatform.ClosePositionView
			view, err = platform.ClosePosition(form)
			if err != nil {
				return
			}
			logger.Infof("Closing position %s for pair %s ask=%0.2f bid=%0.2f make result of %0.4f with amount %0.2f (%0.3f%%)",
				position.ID, position.Pair, view.AskPrice, actualBidPrice, view.Result, position.Amount, resultInPercent)
			var positionToUpdate *repositoryModel.Position
			positionToUpdate, err = svc.positionRepo.Get(position.ID)
			if err != nil {
				return
			}
			positionToUpdate.Closed = true
			positionToUpdate.BidPrice = view.BidPrice
			positionToUpdate.Result = view.Result
			positionToUpdate.BidDate = time.Now()
		}
	}
	return
}

func (svc *service) shouldOpenPosition(actualBalance float64, lastPrices []*repositoryModel.Price, nbOpenedPosition int) (shouldOpen bool, err error) {
	// Don't open if not enough in balance
	if nbOpenedPosition == minimumAmountToOpenPosition {
		shouldOpen = false
		return
	}
	if actualBalance == 0 {
		shouldOpen = false
		return
	}
	if svc.getAmountToInvest(actualBalance, nbOpenedPosition) < minimumAmountToOpenPosition {
		shouldOpen = false
		return
	}

	// don't open if position has been already opened using one of this prices
	shouldOpen = true
	for i := 1; i < nbOfIncreasingValueToOpenPosition+1; i++ {
		previous := lastPrices[i-1].AskPrice
		actual := lastPrices[i].AskPrice
		if actual < previous {
			shouldOpen = false
			return
		}
		_, ok := svc.openedPositionIdsByTimestamp[lastPrices[i].Date.Unix()]
		if ok {
			shouldOpen = false
			return
		}
	}
	return
}

func (svc *service) getAmountToInvest(currentBalance float64, nbOpenedPosition int) (amount float64) {
	if nbOpenedPosition != maxOpenedPositionsByCurrency {
		amount = currentBalance / float64(maxOpenedPositionsByCurrency-nbOpenedPosition)
	}
	return
}
