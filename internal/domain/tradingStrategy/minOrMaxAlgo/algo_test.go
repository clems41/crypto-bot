package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/domain/tradingStrategy"
	"crypto-bot/internal/model"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_SellOrderOpen(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	pair3 := tradingConst.EthereumBitcoin
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2, pair3},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 50,
		},
		OpenOrders: []model.Order{
			{
				Pair:   pair1,
				Side:   tradingConst.Sell,
				Status: tradingConst.Open,
			},
			{
				Pair:   pair3,
				Side:   tradingConst.Buy,
				Status: tradingConst.Open,
			},
		},
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.True(t, ok)

	// check order amount
	require.Equal(t, 50.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_BuyOrderOpen(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	pair3 := tradingConst.EthereumBitcoin
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2, pair3},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 100,
		},
		OpenOrders: []model.Order{
			{
				Pair:   pair1,
				Side:   tradingConst.Buy,
				Status: tradingConst.Open,
			},
		},
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.True(t, ok)

	// check order amount
	require.Equal(t, 50.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_MaximumBuyOrder(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 100,
		},
		OpenOrders: []model.Order{
			{
				Pair:   pair1,
				Side:   tradingConst.Buy,
				Status: tradingConst.Open,
			},
			{
				Pair:   pair2,
				Side:   tradingConst.Buy,
				Status: tradingConst.Open,
			},
		},
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.False(t, ok)

	// check order amount
	require.Equal(t, 0.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_MaximumSellOrder(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 100,
		},
		OpenOrders: []model.Order{
			{
				Pair:   pair1,
				Side:   tradingConst.Buy,
				Status: tradingConst.Open,
			},
			{
				Pair:   pair2,
				Side:   tradingConst.Sell,
				Status: tradingConst.Open,
			},
		},
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.False(t, ok)

	// check order amount
	require.Equal(t, 0.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_NotEnoughInBalance(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	pair3 := tradingConst.EthereumBitcoin
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2, pair3},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 5,
		},
		OpenOrders: []model.Order{
			{
				Pair:   pair3,
				Side:   tradingConst.Sell,
				Status: tradingConst.Open,
			},
		},
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.False(t, ok)

	// check order amount
	require.Equal(t, 0.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_ZeroPreviousOrder(t *testing.T) {
	// instantiate algo
	platformName := fake.Sentences()
	pair1 := tradingConst.BitcoinEuro
	pair2 := tradingConst.EthereumEuro
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
		InitialPairsToTradeByPlatform: map[string][]tradingConst.Pair{
			platformName: {pair1, pair2},
		},
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	askPrice := 5432.2
	form := tradingStrategy.ShouldAddOrderForm{
		PriceHistory: []model.Price{
			{
				Pair: pair2,
				Ask:  askPrice,
			},
		},
		IndexPrice: model.Price{
			Pair: pair2,
			Ask:  askPrice - 100,
		},
		PairToTrade: pair2,
		CurrentBalance: map[tradingConst.Currency]float64{
			tradingConst.Euro: 100,
		},
		OpenOrders:   nil,
		PlatformName: platformName,
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.True(t, ok)

	// check order amount
	require.Equal(t, 50.0, order.Amount)
}
