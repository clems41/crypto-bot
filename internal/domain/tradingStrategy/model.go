package tradingStrategy

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
)

/* Should Add Order */

type ShouldAddOrderForm struct {
	PriceHistory   []model.Price
	IndexPrice     model.Price
	PairToTrade    tradingConst.Pair
	CurrentBalance map[tradingConst.Currency]float64
	OpenOrders     []model.Order
	AllPairsTraded []tradingConst.Pair
}
