package tradingPlatformMock

import (
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
	require.Equal(t, initWalletBalance, view.Balance)
}

func TestMock_GetPrice(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	form := tradingPlatform.GetPriceForm{
		Currency: fake.Word(),
	}
	for range fakeData.FakeRange(25, 50) {
		var view tradingPlatform.GetPriceView
		view, err = testMock.GetPrice(form)
		require.NoError(t, err)
		require.True(t, view.Value > 10000 && view.Value < 70000) // In all dataset, bitcoin must be between 10 000 and 70 000
	}
}

func TestMock_OpenPosition(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	form := tradingPlatform.OpenPositionForm{
		Currency: fake.Word(),
		Amount:   initWalletBalance / 10,
	}
	view, err := testMock.OpenPosition(form)
	require.NoError(t, err)
	require.NotEqual(t, "", view.PositionID)
}

func TestMock_ClosePosition(t *testing.T) {
	testMock, err := New()
	require.NoError(t, err)
	form := tradingPlatform.ClosePositionForm{
		PositionID: fake.Word(),
	}
	view, err := testMock.ClosePosition(form)
	require.Error(t, err) // should return error because position doesn't exist

	// create existing position
	openForm := tradingPlatform.OpenPositionForm{
		Currency: fake.Word(),
		Amount:   initWalletBalance / 10,
	}
	openView, err := testMock.OpenPosition(openForm)
	require.NoError(t, err)
	require.NotEqual(t, "", openView.PositionID)

	// getting some price to get new value from dataset
	for range fakeData.FakeRange(5, 20) {
		_, err = testMock.GetPrice(tradingPlatform.GetPriceForm{
			Currency: fake.Word(),
		})
		require.NoError(t, err)
	}

	// try again with created position
	form = tradingPlatform.ClosePositionForm{
		PositionID: openView.PositionID,
	}
	view, err = testMock.ClosePosition(form)
	require.NoError(t, err) // should not return error because position exists
	require.NotEqual(t, float32(0), view.Result)
	require.NotEqual(t, openForm.Amount, view.Result)
}
