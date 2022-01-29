package tradingStrategy

import "crypto-bot/internal/model"

/* Should Add Order */

type ShouldAddOrderForm struct {
	PriceHistory   []model.Price
	IndexPrice     model.Price
	PairToTrade    string
	CurrentBalance map[string]float64
	OpenOrders     []model.Order
	AllPairsTraded []string
}
