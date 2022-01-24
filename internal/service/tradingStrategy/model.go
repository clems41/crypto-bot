package tradingStrategy

import "crypto-bot/internal/model"

/* Should Add Order */

type ShouldAddOrderForm struct {
	PriceHistory []model.Price
	IndexPrice   model.Price
}

type ShouldAddOrderView struct {
	ShouldAddOrder bool
	Side           string
	Type           string
	Price          float64
}
