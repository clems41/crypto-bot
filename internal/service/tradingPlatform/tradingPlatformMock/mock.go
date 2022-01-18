package tradingPlatformMock

import (
	"crypto-bot/internal/constant/currencyConst"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/testUtils/fakeData"
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
	dataset, err := openDataset2()
	if err != nil {
		return nil, err
	}

	newMock := &mock{
		balanceByCurrency: map[string]float64{
			currencyConst.BtcEurPair: initWalletBalance,
		},
		actualDatasetIndex: initDatasetIndex,
		positions:          make(map[string]tradingPlatform.PositionView),
		openedPositions:    make(map[string]tradingPlatform.PositionView),
		dataset:            dataset,
	}
	return newMock, nil
}

func (mock *mock) GetWalletBalance() (view tradingPlatform.WalletView, err error) {
	return tradingPlatform.WalletView{
		BalanceByCurrency: mock.balanceByCurrency,
	}, nil
}

func (mock *mock) OpenPosition(form tradingPlatform.OpenPositionForm) (view tradingPlatform.OpenPositionView, err error) {
	// create new position
	positionID := fakeData.UuidWithOnlyAlphaNumeric()
	position := tradingPlatform.PositionView{
		ID:       positionID,
		Currency: form.Currency,
		Amount:   form.Amount,
		Price:    form.AskPrice,
	}

	// adding position in fake platform
	mock.positions[positionID] = position
	mock.openedPositions[positionID] = position

	// decrease balance
	mock.balanceByCurrency[form.Currency] -= form.Amount

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
	askPrice := position.Price
	bidPrice := mock.dataset[mock.actualDatasetIndex].BidPrice
	profit := bidPrice * position.Amount / askPrice
	result := profit - position.Amount
	position.Result = result

	// increase balance
	mock.balanceByCurrency[position.Currency] += profit

	// store position with result
	mock.positions[position.ID] = position

	// close position
	delete(mock.openedPositions, form.PositionID)

	// fill view
	view = tradingPlatform.ClosePositionView{
		Result: result,
	}
	return
}

func (mock *mock) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	price := mock.dataset[mock.actualDatasetIndex]
	view = tradingPlatform.GetPriceView{
		AskPrice: price.AskPrice,
		BidPrice: price.BidPrice,
		Date:     price.Date,
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
			ID:     position.ID,
			Amount: position.Amount,
			Price:  position.Price,
			Result: position.Result,
		})
	}
	return
}

func (mock *mock) GetAllPositions() (view tradingPlatform.GetAllPositionsView, err error) {
	for _, position := range mock.positions {
		view.Positions = append(view.Positions, tradingPlatform.PositionView{
			ID:     position.ID,
			Amount: position.Amount,
			Price:  position.Price,
			Result: position.Result,
		})
	}
	return
}
