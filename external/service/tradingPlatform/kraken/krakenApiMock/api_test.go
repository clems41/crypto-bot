package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestApi_GetIndexPrice(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)
	pair := tradingConst.BitcoinEuro
	fakePair := tradingConst.Pair(fake.Word())

	// try with non-existing pair, should return error
	prices, err := apiTest.GetIndexPrices(fakePair)
	require.Error(t, err)

	// try with existing pair, should be ok
	prices, err = apiTest.GetIndexPrices(pair)
	require.NoError(t, err)

	// check view
	require.Len(t, prices, 1)
	require.True(t, prices[0].Bid > 10000 && prices[0].Bid < 70000,
		"Bitcoin bid price should be between 10 000 and 70 000 / unit but it is %0.2f", prices[0].Bid)
	require.True(t, prices[0].Ask > 10000 && prices[0].Ask < 70000,
		"Bitcoin ask price should be between 10 000 and 70 000 / unit but it is %0.2f", prices[0].Ask)
	require.True(t, prices[0].Date.Before(time.Now()) && prices[0].Date.After(time.Now().AddDate(0, 0, -1)),
		"Date price should be between yesterday and now but it is %s", prices[0].Date.String())
}
