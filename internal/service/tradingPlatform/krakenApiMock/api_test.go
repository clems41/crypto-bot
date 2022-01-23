package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
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

func TestApi_GetIndexPrice(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)
	pair := tradingConst.BtcEurPair

	// try with non-existing pair, should return error
	view, err := apiTest.GetIndexPrices(pair)
	require.Error(t, err)

	// try with existing pair, should be ok
	view, err = apiTest.GetIndexPrices(pair)
	require.NoError(t, err)

	// check view
	price, ok := view[pair]
	require.True(t, ok)
	require.True(t, price.Bid > 10000 && price.Bid < 70000,
		"Bitcoin bid price should be between 10 000 and 70 000 / unit but it is %0.2f", price.Bid)
	require.True(t, price.Ask > 10000 && price.Ask < 70000,
		"Bitcoin ask price should be between 10 000 and 70 000 / unit but it is %0.2f", price.Ask)
	require.True(t, price.Date.Before(time.Now()) && price.Date.After(time.Now().AddDate(0, 0, -1)),
		"Date price should be between yesterday and now but it is %s", price.Date.String())
}

func TestApi_GetWalletBalance(t *testing.T) {

}

func TestApi_OpenPosition(t *testing.T) {

}
