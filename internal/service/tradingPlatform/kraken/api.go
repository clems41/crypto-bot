package kraken

import (
	"crypto-bot/internal/constant/currencyConst"
	"crypto-bot/internal/service/tradingPlatform"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var _ tradingPlatform.Api = (*api)(nil)

type api struct {
	balanceByCurrency map[string]float64
	openedPositions   map[string]tradingPlatform.PositionView
	positions         map[string]tradingPlatform.PositionView
}

func New() (*api, error) {
	return &api{
		balanceByCurrency: map[string]float64{
			currencyConst.BtcEurPair: initBalance,
		},
		openedPositions: make(map[string]tradingPlatform.PositionView),
		positions:       make(map[string]tradingPlatform.PositionView),
	}, nil
}

func (api *api) GetOpenedPositions() (view tradingPlatform.GetOpenedPositionsView, err error) {
	for _, position := range api.openedPositions {
		view.Positions = append(view.Positions, position)
	}
	return
}

func (api *api) GetAllPositions() (view tradingPlatform.GetAllPositionsView, err error) {
	for _, position := range api.positions {
		view.Positions = append(view.Positions, position)
	}
	return
}

func (api *api) GetWalletBalance() (view tradingPlatform.WalletView, err error) {
	view = tradingPlatform.WalletView{BalanceByCurrency: api.balanceByCurrency}
	return
}

func (api *api) OpenPosition(form tradingPlatform.OpenPositionForm) (view tradingPlatform.OpenPositionView, err error) {
	positionID := uuid.New().String()
	priceView, err := api.GetPrice(tradingPlatform.GetPriceForm{
		Currency: form.Currency,
	})
	if err != nil {
		return
	}
	position := tradingPlatform.PositionView{
		ID:       positionID,
		Amount:   form.Amount,
		Currency: form.Currency,
		Price:    priceView.AskPrice,
	}
	api.openedPositions[positionID] = position
	api.positions[positionID] = position
	return
}

func (api *api) ClosePosition(form tradingPlatform.ClosePositionForm) (view tradingPlatform.ClosePositionView, err error) {
	// Get bid price
	position, ok := api.openedPositions[form.PositionID]
	if !ok {
		return view, errPositionNotFound
	}
	priceView, err := api.GetPrice(tradingPlatform.GetPriceForm{
		Currency: position.Currency,
	})
	if err != nil {
		return
	}

	// calculate profit
	diff := position.Price - priceView.BidPrice
	profit := diff / position.Price * position.Amount

	// store position result
	position.Result = profit
	api.positions[position.ID] = position

	// close position
	delete(api.openedPositions, form.PositionID)
	return
}

func (api *api) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	// get prices from api
	krakenPair, ok := currencyConversion[form.Currency]
	if !ok {
		return view, errCurrencyNotFound
	}
	krakenApi := krakenapi.New(apiKey, apiSecret)
	result, err := krakenApi.Ticker(krakenPair)
	if err != nil {
		return
	}

	// convert value
	var askPrice, bidPrice float64
	if len(result.GetPairTickerInfo(krakenPair).Ask[0]) > 0 {
		askPriceStr := result.GetPairTickerInfo(krakenPair).Ask[0]
		askPrice, err = strconv.ParseFloat(askPriceStr, 64)
		if err != nil {
			return
		}
	} else {
		return view, errCannotFindPriceFromResponse
	}
	if len(result.GetPairTickerInfo(krakenPair).Bid[0]) > 0 {
		bidPriceStr := result.GetPairTickerInfo(krakenPair).Bid[0]
		bidPrice, err = strconv.ParseFloat(bidPriceStr, 64)
		if err != nil {
			return
		}
	} else {
		return view, errCannotFindPriceFromResponse
	}

	// fill view
	view = tradingPlatform.GetPriceView{
		AskPrice: askPrice,
		BidPrice: bidPrice,
		Date:     time.Now(),
	}

	return
}
