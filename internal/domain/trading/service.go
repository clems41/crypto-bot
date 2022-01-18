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
}

type service struct {
	quitChannel     chan bool           // quit goroutine when program exit
	cryptoAPI       tradingPlatform.Api // communicate with trading platform
	priceRepository repository.Price    // use to store and get previous prices
}

func NewService(cryptoAPI tradingPlatform.Api, priceRepo repository.Price) Service {
	return &service{
		cryptoAPI:       cryptoAPI,
		quitChannel:     make(chan bool),
		priceRepository: priceRepo,
	}
}

// Start will run trading algorithm.
// This method must be called in goroutine that will be stopped by Stop method.
func (svc *service) Start() (err error) {
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

	// Close all positions
	openedPositions, err := svc.cryptoAPI.GetOpenedPositions()
	if err != nil {
		return
	}
	for _, position := range openedPositions.Positions {
		form := tradingPlatform.ClosePositionForm{PositionID: position.ID}
		_, err = svc.cryptoAPI.ClosePosition(form)
		if err != nil {
			return
		}
	}

	// Get all positions
	positions, err := svc.cryptoAPI.GetAllPositions()
	if err != nil {
		return
	}
	for _, position := range positions.Positions {
		logger.Infof("Position %s ask=%0.2f bid=%0.2f make result of %0.4f with amount %0.2f (%0.3f%%)",
			position.ID, position.AskPrice, position.BidPrice, position.Result, position.Amount,
			tradingUtils.GetResultInPercent(position.Amount, position.Result))
	}

	// Update current balance
	balanceView, err := svc.cryptoAPI.GetWalletBalance()
	if err != nil {
		return
	}
	logger.Infof("Current balance is now of %v", balanceView.BalanceByCurrency)

	return
}

func (svc *service) applyTradingAlgorithm() (err error) {
	// Update current balance
	balanceView, err := svc.cryptoAPI.GetWalletBalance()
	if err != nil {
		return
	}
	logger.Infof("Current balance is %v", balanceView.BalanceByCurrency)

	// Loop over all currencies to trade
	for _, currency := range currenciesToTrade {
		// Getting actual price of currency
		priceForm := tradingPlatform.GetPriceForm{
			Currency: currency,
		}
		var priceView tradingPlatform.GetPriceView
		priceView, err = svc.cryptoAPI.GetPrice(priceForm)
		if err != nil {
			return
		}
		logger.Infof("%s price is ask=%0.2f and bid=%0.2f", currency, priceView.AskPrice, priceView.BidPrice)

		// Store new price
		price := repositoryModel.Price{
			Date:     priceView.Date,
			Currency: currency,
			AskPrice: priceView.AskPrice,
			BidPrice: priceView.BidPrice,
		}
		err = svc.priceRepository.Store(&price)

		// Get last prices
		var lastPrices []*repositoryModel.Price
		lastPrices, err = svc.priceRepository.GetLast(nbPreviousValues, currency)
		if err != nil {
			return
		}

		// Get opened positions
		var openedPositions tradingPlatform.GetOpenedPositionsView
		openedPositions, err = svc.cryptoAPI.GetOpenedPositions()
		if err != nil {
			return
		}

		// Open and close positions depending on new prices and previous positions
		if len(openedPositions.Positions) >= maxOpenedPositions {
			logger.Infof("No new position will be opened because max %d has been reached", maxOpenedPositions)
		} else {
			err = svc.openNewPositions(balanceView.BalanceByCurrency[currency], lastPrices, currency)
			if err != nil {
				return
			}
		}

		err = svc.closePositions(priceView.BidPrice, openedPositions.Positions)
		if err != nil {
			return
		}
	}

	return
}

func (svc *service) openNewPositions(actualBalance float64, lastPrices []*repositoryModel.Price, currency string) (err error) {
	// if nbOfIncreasingValueToOpenPosition last values are all increasing, we should open new position
	shouldOpen := true
	for i := nbPreviousValues - nbOfIncreasingValueToOpenPosition; i < nbPreviousValues; i++ {
		previous := lastPrices[i-1].AskPrice
		actual := lastPrices[i].AskPrice
		if previous > actual {
			shouldOpen = false
		}
	}

	if shouldOpen {
		lastPrice := lastPrices[len(lastPrices)-1]
		form := tradingPlatform.OpenPositionForm{
			AskPrice: lastPrice.AskPrice,
			Currency: currency,
			Amount:   actualBalance * 3 / 4, // use only 3/4 of available balance
		}
		var view tradingPlatform.OpenPositionView
		view, err = svc.cryptoAPI.OpenPosition(form)
		logger.Infof("Opening new position %s with amount of %f", view.PositionID, form.Amount)
		if err != nil {
			return
		}
	}
	return
}

func (svc *service) closePositions(bidPrice float64, openedPositions []tradingPlatform.PositionView) (err error) {
	// close position if actual price is more than 0.1% of price position
	for _, position := range openedPositions {
		if bidPrice >= (1+gainToClosePositionInPercent/100)*position.AskPrice {
			form := tradingPlatform.ClosePositionForm{PositionID: position.ID}
			var view tradingPlatform.ClosePositionView
			view, err = svc.cryptoAPI.ClosePosition(form)
			if err != nil {
				return
			}
			logger.Infof("Closing position %s ask=%0.2f bid=%0.2f make result of %0.4f with amount %0.2f (%0.3f%%)",
				position.ID, position.AskPrice, bidPrice, view.Result, position.Amount,
				tradingUtils.GetResultInPercent(position.Amount, view.Result))
		}
	}
	return
}
