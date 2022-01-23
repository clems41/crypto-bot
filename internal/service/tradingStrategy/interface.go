package tradingStrategy

import "crypto-bot/internal/model"

type Algo interface {
	ShouldAddOrder(form ShouldAddOrderForm) (ok bool, order *model.Order, err error)
}
