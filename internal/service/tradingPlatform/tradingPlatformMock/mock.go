package tradingPlatformMock

import (
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"fmt"
)

var _ tradingPlatform.Api = (*mock)(nil)

type mock struct {
	balance            float64
	actualDatasetIndex int
	positions          map[string]Position
	openedPositions    map[string]Position
	dataset            []float64
}

func New() (*mock, error) {
	dataset, err := openDataset2()
	if err != nil {
		return nil, err
	}

	newMock := &mock{
		balance:            initWalletBalance,
		actualDatasetIndex: initDatasetIndex,
		positions:          make(map[string]Position),
		openedPositions:    make(map[string]Position),
		dataset:            dataset,
	}
	return newMock, nil
}

func (mock *mock) GetWalletBalance() (view tradingPlatform.WalletView, err error) {
	return tradingPlatform.WalletView{
		Balance: mock.balance,
	}, nil
}

func (mock *mock) OpenPosition(form tradingPlatform.OpenPositionForm) (view tradingPlatform.OpenPositionView, err error) {
	positionID := fakeData.UuidWithOnlyAlphaNumeric()
	position := Position{
		ID:           positionID,
		DatasetIndex: mock.actualDatasetIndex,
		Amount:       form.Amount,
	}
	mock.positions[positionID] = position
	mock.openedPositions[positionID] = position
	mock.balance -= form.Amount
	view = tradingPlatform.OpenPositionView{PositionID: positionID}
	return
}

func (mock *mock) ClosePosition(form tradingPlatform.ClosePositionForm) (view tradingPlatform.ClosePositionView, err error) {
	position, ok := mock.positions[form.PositionID]
	if !ok {
		return view, fmt.Errorf("position %s not found", form.PositionID)
	}
	startValue := mock.dataset[position.DatasetIndex]
	stopValue := mock.dataset[mock.actualDatasetIndex]
	coefficient := (stopValue-startValue)/startValue + 1
	result := position.Amount * coefficient
	mock.balance += result
	position.Result = result
	mock.positions[position.ID] = position
	delete(mock.openedPositions, form.PositionID)
	view = tradingPlatform.ClosePositionView{
		Result: result,
	}
	return
}

func (mock *mock) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	value := mock.dataset[mock.actualDatasetIndex]
	view = tradingPlatform.GetPriceView{
		Value: float64(value),
	}
	mock.actualDatasetIndex++
	if mock.actualDatasetIndex >= len(mock.dataset) {
		mock.actualDatasetIndex = 0
	}
	return
}

func (mock *mock) GetOpenedPositions() (view tradingPlatform.GetOpenedPositionsView, err error) {
	for _, position := range mock.openedPositions {
		value := mock.dataset[position.DatasetIndex]
		view.Positions = append(view.Positions, tradingPlatform.PositionView{
			ID:     position.ID,
			Amount: position.Amount,
			Price:  value,
			Result: position.Result,
		})
	}
	return
}

func (mock *mock) GetAllPositions() (view tradingPlatform.GetAllPositionsView, err error) {
	for _, position := range mock.positions {
		value := mock.dataset[position.DatasetIndex]
		view.Positions = append(view.Positions, tradingPlatform.PositionView{
			ID:     position.ID,
			Amount: position.Amount,
			Price:  value,
			Result: position.Result,
		})
	}
	return
}
