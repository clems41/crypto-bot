package tradingStrategy

import "crypto-bot/internal/model"

type Algo interface {
	// ShouldAddOrder determines if order should be open. If returned order is nil, order should not be open
	ShouldAddOrder(form ShouldAddOrderForm) (shouldOpen bool, order model.Order, err error)
	// PricesNeeded return number of prices needed by ShouldAddOrder to work.
	// If ShouldAddOrder got fewer prices than value returned by PricesNeeded, ShouldAddOrder will return an error.
	PricesNeeded() (numberOfPrices int)
}
