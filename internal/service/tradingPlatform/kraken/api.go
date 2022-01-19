package kraken

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/envUtils"
	"crypto-bot/pkg/utils/tradingUtils"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var _ tradingPlatform.Api = (*api)(nil)

type api struct {
	apiKey            string
	apiSecret         string
	balanceByCurrency map[string]float64
	openedPositions   map[string]tradingPlatform.PositionView
	positions         map[string]tradingPlatform.PositionView
}

func New() (*api, error) {
	apiKey, err := envUtils.GetFromEnvOrError(envKrakenApiKey)
	if err != nil {
		return nil, err
	}
	apiSecret, err := envUtils.GetFromEnvOrError(envKrakenApiSecret)
	if err != nil {
		return nil, err
	}

	return &api{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		balanceByCurrency: map[string]float64{
			tradingConst.EuroCurrency: initBalance,
		},
		openedPositions: make(map[string]tradingPlatform.PositionView),
		positions:       make(map[string]tradingPlatform.PositionView),
	}, nil
}

func (api *api) Name() (name string) {
	return tradingConst.KrakenPlatform
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
	// Check balance
	currency := tradingConst.CurrencyNeededToTradePair(form.Pair)
	balance, ok := api.balanceByCurrency[currency]
	if !ok {
		return view, tradingPlatform.ErrCurrencyNotInWallet
	}
	if form.Amount > balance {
		return view, tradingPlatform.ErrNotEnoughCash
	}

	// Create position
	positionID := uuid.New().String()
	priceView, err := api.GetPrice(tradingPlatform.GetPriceForm{
		Pairs: []string{form.Pair},
	})
	if err != nil {
		return
	}
	price, ok := priceView.PriceByPair[form.Pair]
	if !ok {
		return view, tradingPlatform.ErrPairNotFound
	}
	position := tradingPlatform.PositionView{
		ID:       positionID,
		Amount:   form.Amount * (1 - fakeFeesInPercent/100), // trading platform always keep little percent of invest
		Pair:     form.Pair,
		AskPrice: price.AskPrice,
	}

	// Push order
	api.openedPositions[positionID] = position
	api.positions[positionID] = position

	// Update wallet balance
	api.balanceByCurrency[currency] = balance - form.Amount
	view = tradingPlatform.OpenPositionView{PositionID: positionID}
	return
}

func (api *api) ClosePosition(form tradingPlatform.ClosePositionForm) (view tradingPlatform.ClosePositionView, err error) {
	// Get bid price
	position, ok := api.openedPositions[form.PositionID]
	if !ok {
		return view, tradingPlatform.ErrPositionNotFound
	}
	priceView, err := api.GetPrice(tradingPlatform.GetPriceForm{
		Pairs: []string{position.Pair},
	})
	if err != nil {
		return
	}

	// calculate profit
	price, ok := priceView.PriceByPair[position.Pair]
	if !ok {
		return view, tradingPlatform.ErrPairNotFound
	}
	askPrice := position.AskPrice
	bidPrice := price.BidPrice
	profit := tradingUtils.GetProfit(askPrice, bidPrice, position.Amount)
	result := tradingUtils.GetResult(askPrice, bidPrice, position.Amount)
	api.balanceByCurrency[position.Pair] += profit

	// store position result
	position.Result = result
	position.BidPrice = bidPrice
	api.positions[position.ID] = position

	// close position
	delete(api.openedPositions, form.PositionID)
	view = tradingPlatform.ClosePositionView{
		AskPrice: askPrice,
		BidPrice: bidPrice,
		Result:   result,
		Profit:   profit,
	}
	return
}

func (api *api) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	// convert pairs
	var krakenPairs []string
	for _, pair := range form.Pairs {
		var krakenPair string
		krakenPair, err = GetKrakenPair(pair)
		if err != nil {
			return
		}
		krakenPairs = append(krakenPairs, krakenPair)
	}

	// get prices from api
	krakenApi := krakenapi.New(api.apiKey, api.apiSecret)
	result, err := krakenApi.Ticker(krakenPairs...)
	if err != nil {
		return
	}
	if result == nil {
		return view, tradingPlatform.ErrEmptyResponse
	}

	// convert prices from response
	view.PriceByPair = make(map[string]tradingPlatform.PriceView)
	for _, krakenPair := range krakenPairs {
		var askPrice, bidPrice float64
		if len(result.GetPairTickerInfo(krakenPair).Ask) > 0 {
			askPriceStr := result.GetPairTickerInfo(krakenPair).Ask[0]
			askPrice, err = strconv.ParseFloat(askPriceStr, 64)
			if err != nil {
				return
			}
		} else {
			return view, tradingPlatform.ErrCannotFindPriceFromResponse
		}
		if len(result.GetPairTickerInfo(krakenPair).Bid[0]) > 0 {
			bidPriceStr := result.GetPairTickerInfo(krakenPair).Bid[0]
			bidPrice, err = strconv.ParseFloat(bidPriceStr, 64)
			if err != nil {
				return
			}
		} else {
			return view, tradingPlatform.ErrCannotFindPriceFromResponse
		}

		// fill view
		price := tradingPlatform.PriceView{
			AskPrice: askPrice,
			BidPrice: bidPrice,
			Date:     time.Now(),
		}
		var pair string
		pair, err = GetProjectPair(krakenPair)
		if err != nil {
			return
		}
		view.PriceByPair[pair] = price
	}

	return
}
