package tradingPlatform

import "crypto-bot/internal/model"

type Api interface {
	Name() (name string)
	AddOrder(order *model.Order) (err error)
	CancelOrder(orderID string) (err error)
	CancelAllOrders() (view CancelAllOrdersView, err error)
	GetPrices(form GetPriceForm) (prices []*model.Price, err error)
	GetIndexPrices(pairs ...string) (prices map[string]*model.Price, err error)
	GetOpenOrders() (orders []*model.Order, err error)
	GetAllOrders() (orders []*model.Order, err error)
	GetBalance() (balance *model.Balance, err error)
}
