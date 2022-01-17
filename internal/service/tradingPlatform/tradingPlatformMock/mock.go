package tradingPlatformMock

import (
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

var _ tradingPlatform.Api = (*mock)(nil)

type mock struct {
	balance            float32
	actualDatasetIndex int
	positions          map[string]Position
	dataset            Dataset
}

func New() (*mock, error) {
	datasetFile, err := ioutil.ReadFile(datasetPath)
	if err != nil {
		return nil, err
	}
	var dataset Dataset
	err = json.Unmarshal(datasetFile, &dataset)
	if err != nil {
		return nil, err
	}
	newMock := &mock{
		balance:            initWalletBalance,
		actualDatasetIndex: 0,
		positions:          make(map[string]Position),
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
	mock.positions[positionID] = Position{
		ID:           positionID,
		DatasetIndex: mock.actualDatasetIndex,
		Amount:       form.Amount,
	}
	view = tradingPlatform.OpenPositionView{PositionID: positionID}
	return
}

func (mock *mock) ClosePosition(form tradingPlatform.ClosePositionForm) (view tradingPlatform.ClosePositionView, err error) {
	position, ok := mock.positions[form.PositionID]
	if !ok {
		return view, fmt.Errorf("position %s not found", form.PositionID)
	}
	startValueStr := mock.dataset.Values[position.DatasetIndex].Hp
	stopValueStr := mock.dataset.Values[mock.actualDatasetIndex].Hp
	startValue, err := strconv.ParseFloat(startValueStr, 64)
	if err != nil {
		return
	}
	stopValue, err := strconv.ParseFloat(stopValueStr, 64)
	if err != nil {
		return
	}
	coefficient := float32((stopValue-startValue)/stopValue + 1)
	result := position.Amount * coefficient
	view = tradingPlatform.ClosePositionView{
		Result: result,
	}
	return
}

func (mock *mock) GetPrice(form tradingPlatform.GetPriceForm) (view tradingPlatform.GetPriceView, err error) {
	valueStr := mock.dataset.Values[mock.actualDatasetIndex].Hp
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return
	}
	view = tradingPlatform.GetPriceView{
		Value: float32(value),
	}
	mock.actualDatasetIndex++
	if mock.actualDatasetIndex >= len(mock.dataset.Values) {
		mock.actualDatasetIndex = 0
	}
	return
}
