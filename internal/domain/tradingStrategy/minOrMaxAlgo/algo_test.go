package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/domain/tradingStrategy"
	"crypto-bot/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_OneOrderOpen(t *testing.T) {
	// instantiate algo
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	pair1 := tradingConst.BtcEurPair
	pair2 := tradingConst.EthEurPair
	pair3 := tradingConst.EthBtcPair
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
		CurrentBalance: map[string]float64{
			tradingConst.EuroCurrency: 50,
		},
		NbOpenOrdersByPair: map[string]int{
			pair1: 1,
			pair2: 0,
			pair3: 1,
		},
		AllPairsTraded: []string{pair1, pair2, pair3},
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.True(t, ok)

	// check order amount
	require.Equal(t, 50.0, order.Amount)

}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_MaximumOrderByPair(t *testing.T) {
	// instantiate algo
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	pair1 := tradingConst.BtcEurPair
	pair2 := tradingConst.EthEurPair
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
		CurrentBalance: map[string]float64{
			tradingConst.EuroCurrency: 100,
		},
		NbOpenOrdersByPair: map[string]int{
			pair1: 0,
			pair2: 1,
		},
		AllPairsTraded: []string{pair1, pair2},
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.False(t, ok)

	// check order amount
	require.Equal(t, 0.0, order.Amount)

}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_NotEnoughInBalance(t *testing.T) {
	// instantiate algo
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	pair1 := tradingConst.BtcEurPair
	pair2 := tradingConst.EthEurPair
	pair3 := tradingConst.EthBtcPair
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
		CurrentBalance: map[string]float64{
			tradingConst.EuroCurrency: 5,
		},
		NbOpenOrdersByPair: map[string]int{
			pair1: 0,
			pair2: 0,
			pair3: 1,
		},
		AllPairsTraded: []string{pair1, pair2, pair3},
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.False(t, ok)

	// check order amount
	require.Equal(t, 0.0, order.Amount)
}

func TestAlgorithm_ShouldAddOrder_GetAmountToBuy_ZeroPreviousOrder(t *testing.T) {
	// instantiate algo
	config := Config{
		NumberOfPreviousPricesToCompare:       1,
		PercentPriceBelowToBuy:                0,
		MinimumResultInPercentToClosePosition: 0.5,
		MaxOpenedOrdersByPair:                 1,
		MinimumAmount:                         10,
	}
	algo, err := New(&config)
	require.NoError(t, err)

	// get order to open
	pair1 := tradingConst.BtcEurPair
	pair2 := tradingConst.EthEurPair
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
		CurrentBalance: map[string]float64{
			tradingConst.EuroCurrency: 100,
		},
		NbOpenOrdersByPair: map[string]int{
			pair1: 0,
			pair2: 0,
		},
		AllPairsTraded: []string{pair1, pair2},
	}
	ok, order, err := algo.ShouldAddOrder(form)
	require.NoError(t, err)
	require.True(t, ok)

	// check order amount
	require.Equal(t, 50.0, order.Amount)
}
