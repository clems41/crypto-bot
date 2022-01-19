package tradingPlatformMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"crypto-bot/pkg/utils/tradingUtils"
	"fmt"
)

var _ tradingPlatform.Api = (*mock)(nil)

type mock struct {
	balanceByCurrency  map[string]float64
	actualDatasetIndex int
	positions          map[string]tradingPlatform.PositionView
	openedPositions    map[string]tradingPlatform.PositionView
	dataset            []repositoryModel.Price
}

func New() (*mock, error) {
	openDatasetMethod, ok := openDataset[datasetIDToUse]
	if !ok {
		return nil, errDatasetNotFound
	}
	dataset, err := openDatasetMethod()
	if err != nil {
		return nil, err
	}

	newMock := &mock{
		balanceByCurrency: map[string]float64{
			tradingConst.EuroCurrency: initWalletBalance,
		},
		actualDatasetIndex: initDatasetIndex,
		positions:          make(map[string]tradingPlatform.PositionView),
		openedPositions:    make(map[string]tradingPlatform.PositionView),
		dataset:            dataset,
	}
	return newMock, nil
}

func (mock *mock) Name() (name string) {
	return tradingConst.MockPlatform
}

func (mock *mock) GetWalletBalance() (view tradingPlatform.WalletView, err error) {
	return tradingPlatform.WalletView{
		BalanceByCurrency: mock.balanceByCurrency,
	}, nil
}

func (mock *mock) OpenPosition(form tradingPlatform.OpenPositionForm) (view tradingPlatform.OpenPositionView, err error) {
	// check balance
	currency := tradingConst.CurrencyNeededToTradePair(form.Pair)
	balance, ok := mock.balanceByCurrency[currency]
	if !ok {
		return view, tradingPlatform.ErrPairNotFound
	}
	if form.Amount > balance {
		return view, tradingPlatform.ErrNotEnoughCash
	}
	// create new position
	positionID := fakeData.UuidWithOnlyAlphaNumeric()
	position := tradingPlatform.PositionView{
		ID:       positionID,
		Pair:     form.Pair,
		Amount:   form.Amount,
		AskPrice: form.AskPrice,
	}

	// adding position in fake platform
	mock.positions[positionID] = position
	mock.openedPositions[positionID] = position

	// decrease balance
	mock.balanceByCurrency[currency] = balance - form.Amount

	// fill view
	view = tradingPlatform.OpenPositionView{PositionID: positionID}
	return
}

func (mock *mock) ClosePosition(form tradingPlatform.ClosePositionForm) (view tradingPlatform.ClosePositionView, err error) {
	// find position
	position, ok := mock.positions[form.PositionID]
	if !ok {
		return view, fmt.Errorf("position %s not found", form.PositionID)
	}

	// calculate profit
	askPrice := position.AskPrice
	bidPrice := mock.dataset[mock.actualDatasetIndex].BidPrice
	profit := tradingUtils.GetProfit(askPrice, bidPrice, position.Amount)
	result := tradingUtils.GetResult(askPrice, bidPrice, position.Amount)
	position.Result = result
	position.BidPrice = bidPrice

	// increase balance
	mock.balanceByCurrency[position.Pair] += profit

	// store position with result
	mock.positions[position.ID] = position

	// close position
	delete(mock.openedPositions, form.PositionID)

	// fill view
	view = tradingPlatform.ClosePositionView{
		AskPrice: askPrice,
		BidPrice: bidPrice,
		Result:   result,
		Profit:   profit,
	}
	return
}

func (mock *mock) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	price := mock.dataset[mock.actualDatasetIndex]
	view.PriceByPair = make(map[string]tradingPlatform.PriceView)
	for _, pair := range form.Pairs {
		view.PriceByPair[pair] = tradingPlatform.PriceView{
			AskPrice: price.AskPrice,
			BidPrice: price.BidPrice,
			Date:     price.Date,
		}
	}
	mock.actualDatasetIndex++
	if mock.actualDatasetIndex >= len(mock.dataset) {
		mock.actualDatasetIndex = 0
	}
	return
}

func (mock *mock) GetOpenedPositions() (view tradingPlatform.GetOpenedPositionsView, err error) {
	for _, position := range mock.openedPositions {
		view.Positions = append(view.Positions, tradingPlatform.PositionView{
			ID:       position.ID,
			Pair:     position.Pair,
			Amount:   position.Amount,
			AskPrice: position.AskPrice,
			BidPrice: position.BidPrice,
			Result:   position.Result,
		})
	}
	return
}

func (mock *mock) GetAllPositions() (view tradingPlatform.GetAllPositionsView, err error) {
	for _, position := range mock.positions {
		view.Positions = append(view.Positions, tradingPlatform.PositionView{
			ID:       position.ID,
			Pair:     position.Pair,
			Amount:   position.Amount,
			AskPrice: position.AskPrice,
			BidPrice: position.BidPrice,
			Result:   position.Result,
		})
	}
	return
}
