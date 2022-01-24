package tradingPlatform

import "crypto-bot/internal/model"

type Api interface {
	Name() (name string)
	AddOrder(order *model.Order) (err error)
	CancelAllOrders() (view CancelAllOrdersView, err error)
	GetIndexPrices(pairs ...string) (prices []model.Price, err error)
	GetOpenOrders() (orders []model.Order, err error)
	GetAllOrders() (orders []model.Order, err error)
	RefreshBalance(balance *model.Balance) (err error)
}
