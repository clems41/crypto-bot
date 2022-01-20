package tradingPlatformMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMock_GetWalletBalance(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	view, err := testMock.GetWalletBalance()
	require.NoError(t, err)
	for _, balance := range view.BalanceByCurrency {
		require.Equal(t, initWalletBalance, balance)
	}
}

func TestMock_GetPrice(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	pair := fake.Word()
	form := tradingPlatform.GetPriceForm{
		Pairs: []string{pair},
	}
	for range fakeData.FakeRange(25, 50) {
		var view tradingPlatform.GetPriceView
		view, err = testMock.GetPrice(form)
		require.NoError(t, err)
		price, ok := view.PriceByPair[pair]
		require.True(t, ok)
		require.True(t, price.AskPrice > 10000 && price.AskPrice < 70000) // In all dataset, bitcoin must be between 10 000 and 70 000
		require.True(t, price.BidPrice > 10000 && price.BidPrice < 70000) // In all dataset, bitcoin must be between 10 000 and 70 000
	}
}

func TestMock_OpenPosition(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	amount := initWalletBalance / 10
	pair := tradingConst.BtcEurPair
	currency := tradingConst.EuroCurrency

	// try with not enough money, should return error
	form := tradingPlatform.OpenPositionForm{
		Pair:   pair,
		Amount: 200,
	}
	view, err := testMock.OpenPosition(form)
	require.Error(t, err)

	// try with enough money, should be ok
	form = tradingPlatform.OpenPositionForm{
		Pair:   pair,
		Amount: amount,
	}
	view, err = testMock.OpenPosition(form)
	require.NoError(t, err)
	require.NotEqual(t, "", view.PositionID)

	// balance should be decreased
	balanceView, err := testMock.GetWalletBalance()
	require.NoError(t, err)
	balance, ok := balanceView.BalanceByCurrency[currency]
	require.True(t, ok)
	require.Equal(t, initWalletBalance-amount, balance)

	// position should be opened
	openedPositions, err := testMock.GetOpenedPositions()
	require.NoError(t, err)
	require.Len(t, openedPositions.Positions, 1)
	positions, err := testMock.GetAllPositions()
	require.NoError(t, err)
	require.Len(t, positions.Positions, 1)
}

func TestMock_ClosePosition(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	pair := tradingConst.BtcEurPair
	currency := tradingConst.EuroCurrency
	form := tradingPlatform.ClosePositionForm{
		PositionID: fake.Word(),
	}
	view, err := testMock.ClosePosition(form)
	require.Error(t, err) // should return error because position doesn't exist

	// create existing position
	openForm := tradingPlatform.OpenPositionForm{
		Pair:   pair,
		Amount: initWalletBalance / 10,
	}
	openView, err := testMock.OpenPosition(openForm)
	require.NoError(t, err)
	require.NotEqual(t, "", openView.PositionID)

	// getting some price to get new value from dataset
	for range fakeData.FakeRange(5, 20) {
		_, err = testMock.GetPrice(tradingPlatform.GetPriceForm{
			Pairs: []string{pair},
		})
		require.NoError(t, err)
	}

	// try again with created position
	form = tradingPlatform.ClosePositionForm{
		PositionID: openView.PositionID,
	}
	view, err = testMock.ClosePosition(form)
	require.NoError(t, err) // should not return error because position exists
	require.NotEqual(t, float64(0), view.BidPrice)

	// balance should not the same as start
	wallet, err := testMock.GetWalletBalance()
	require.NoError(t, err)
	balance, ok := wallet.BalanceByCurrency[currency]
	require.True(t, ok)
	require.NotEqual(t, initWalletBalance, balance)

	// position should not be opened
	openedPositions, err := testMock.GetOpenedPositions()
	require.NoError(t, err)
	require.Len(t, openedPositions.Positions, 0)
	positions, err := testMock.GetAllPositions()
	require.NoError(t, err)
	require.Len(t, positions.Positions, 1)
}

func TestMock_GetOpenedPositions(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	pair := tradingConst.BtcEurPair
	view, err := testMock.GetOpenedPositions()
	require.NoError(t, err)
	require.Len(t, view.Positions, 0)

	// try to open positions and check returned values
	var positions []tradingPlatform.OpenPositionView
	for range fakeData.FakeRange(5, 20) {
		var position tradingPlatform.OpenPositionView
		position, err = testMock.OpenPosition(tradingPlatform.OpenPositionForm{
			Pair:   pair,
			Amount: 5,
		})
		require.NoError(t, err)
		positions = append(positions, position)
	}

	// try again with created positions
	view, err = testMock.GetOpenedPositions()
	require.NoError(t, err)
	require.Len(t, view.Positions, len(positions))
}

func TestMock_GetAllPositions(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	pair := tradingConst.BtcEurPair
	view, err := testMock.GetAllPositions()
	require.NoError(t, err)
	require.Len(t, view.Positions, 0)

	// try to open and close positions and check returned values
	var positions []tradingPlatform.OpenPositionView
	for range fakeData.FakeRange(5, 20) {
		var position tradingPlatform.OpenPositionView
		position, err = testMock.OpenPosition(tradingPlatform.OpenPositionForm{
			Pair:   pair,
			Amount: 5,
		})
		require.NoError(t, err)
		positions = append(positions, position)
		_, err = testMock.ClosePosition(tradingPlatform.ClosePositionForm{PositionID: position.PositionID})
		require.NoError(t, err)
	}

	// try again with created positions
	view, err = testMock.GetAllPositions()
	require.NoError(t, err)
	require.Len(t, view.Positions, len(positions))
}
