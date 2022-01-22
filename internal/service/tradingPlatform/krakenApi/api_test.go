package krakenApi

import (
	"crypto-bot/internal/constant/tradingConst"
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
	pair := tradingConst.BtcEurPair

	// try with non-existing pair, should return error
	form := tradingPlatform.GetPriceForm{
		Pairs: []string{fake.Word()},
	}
	view, err := apiTest.GetPrice(form)
	require.Error(t, err)

	// try with existing pair, should be ok
	form = tradingPlatform.GetPriceForm{
		Pairs: []string{pair},
	}
	view, err = apiTest.GetPrice(form)
	require.NoError(t, err)

	// check view
	price, ok := view.PriceByPair[pair]
	require.True(t, ok)
	require.True(t, price.BidPrice > 10000 && price.BidPrice < 70000,
		"Bitcoin bid price should be between 10 000 and 70 000 / unit but it is %0.2f", price.BidPrice)
	require.True(t, price.AskPrice > 10000 && price.AskPrice < 70000,
		"Bitcoin ask price should be between 10 000 and 70 000 / unit but it is %0.2f", price.AskPrice)
	require.True(t, price.Date.Before(time.Now()) && price.Date.After(time.Now().AddDate(0, 0, -1)),
		"Date price should be between yesterday and now but it is %s", price.Date.String())
}

func TestApi_GetWalletBalance(t *testing.T) {

}

func TestApi_OpenPosition(t *testing.T) {

}
