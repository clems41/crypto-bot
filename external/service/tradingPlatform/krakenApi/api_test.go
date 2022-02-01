package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"fmt"
	krakenClient "github.com/beldur/kraken-go-api-client"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestApi_GetIndexPrice(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)
	pair := tradingConst.BtcEurPair

	// try with non-existing pair, should return error
	prices, err := apiTest.GetIndexPrices(fake.Word())
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

func TestKrakenApi_GetOpenOrders(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)

	orders, err := apiTest.GetOpenOrders()
	require.NoError(t, err)
	require.Len(t, orders, 0)
}

func TestKrakenApi_GetAllOrders(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)

	orders, err := apiTest.GetAllOrders(time.Now().AddDate(-1, 0, 0))
	require.NoError(t, err)
	require.True(t, len(orders) > 0, "More than 0 orders should be returned")
}

func TestKrakenApi_convertOrderFromPlatformToProject(t *testing.T) {
	apiTest, err := New()
	require.NoError(t, err)
	openTime := time.Now().AddDate(0, 0, -1)
	closeTime := time.Now()
	pair := tradingConst.EthEurPair
	closePrice := 5423.5
	closeCondition := tradingConst.LimitCloseConditionType
	orderType := tradingConst.MarketOrderType
	price := 5423.1
	side := tradingConst.BuySideOrder
	volume := 0.02548
	fees := 0.26
	status := tradingConst.CloseOrderStatus
	krakenOrder := krakenClient.Order{
		TransactionID: uuid.New().String(),
		Status:        statusConverter[status],
		OpenTime:      float64(openTime.Unix()),
		CloseTime:     float64(closeTime.Unix()),
		Description: krakenClient.OrderDescription{
			AssetPair: "ETHEUR",
			Close:     fmt.Sprintf("close position @ %s %f", typeConverter[closeCondition], closePrice),
			Leverage:  "0",
			OrderType: typeConverter[orderType],
			Type:      sideConverter[side],
		},
		Volume: fmt.Sprintf("%f", volume),
		Fee:    fees,
		Price:  price,
	}

	order, err := apiTest.convertOrderFromPlatformToProject(krakenOrder)
	require.NoError(t, err)
	require.Equal(t, krakenOrder.TransactionID, order.ID)
	require.Equal(t, openTime.Unix(), order.OpenTime.Unix())
	require.Equal(t, closeTime.Unix(), order.CloseTime.Unix())
	require.Equal(t, pair, order.Pair)
	require.Equal(t, side, order.Side)
	require.Equal(t, volume, order.Volume)
	require.Equal(t, orderType, order.Type)
	require.Equal(t, price, order.Price)
	require.Equal(t, order.Price*order.Volume, order.Amount)
	require.Equal(t, 0, order.Leverage)
	require.Equal(t, closeCondition, order.CloseConditionType)
	require.Equal(t, closePrice, order.CloseConditionPrice)
	require.Equal(t, fees, order.Fees)
	require.Equal(t, status, order.Status)
}
