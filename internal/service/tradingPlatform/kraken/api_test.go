package kraken

import (
	"crypto-bot/internal/constant/currencyConst"
	"crypto-bot/internal/service/tradingPlatform"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestApi_ClosePosition(t *testing.T) {

}

func TestApi_GetAllPositions(t *testing.T) {

}

func TestApi_GetOpenedPositions(t *testing.T) {

}

func TestApi_GetPrice(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)

	// try with non-existing currency, should return error
	form := tradingPlatform.GetPriceForm{
		Currency: fake.Word(),
	}
	view, err := apiTest.GetPrice(form)
	require.Error(t, err)

	// try with existing currency, should be ok
	form = tradingPlatform.GetPriceForm{
		Currency: currencyConst.BtcEurPair,
	}
	view, err = apiTest.GetPrice(form)
	require.NoError(t, err)

	// check view
	require.True(t, view.BidPrice > 10000 && view.BidPrice < 70000,
		"Bitcoin bid price should be between 10 000 and 70 000 / unit but it is %0.2f", view.BidPrice)
	require.True(t, view.AskPrice > 10000 && view.AskPrice < 70000,
		"Bitcoin ask price should be between 10 000 and 70 000 / unit but it is %0.2f", view.AskPrice)
	require.True(t, view.Date.Before(time.Now()) && view.Date.After(time.Now().AddDate(0, 0, -1)),
		"Date price should be between yesterday and now but it is %s", view.Date.String())
}

func TestApi_GetWalletBalance(t *testing.T) {

}

func TestApi_OpenPosition(t *testing.T) {

}
