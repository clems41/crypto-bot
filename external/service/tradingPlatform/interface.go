package tradingPlatform

import (
	"crypto-bot/internal/model"
	"time"
)

type Api interface {
	Name() (name string)
	AddOrder(order *model.Order) (err error)
	GetIndexPrices(pairs ...string) (prices []model.Price, err error)
	GetOpenOrders() (orders []model.Order, err error)
	GetAllOrders(since time.Time) (orders []model.Order, err error)
	RefreshBalance(balance *model.Balance) (err error)
}
