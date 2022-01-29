package repository

import "crypto-bot/internal/model"

type Repository interface {
	/* Price */

	StorePrice(price *model.Price) (err error)
	GetPriceHistory(form GetPriceHistoryForm) (prices []model.Price, err error)

	/* Balance */

	StoreBalance(balance *model.Balance) (err error)

	/* order */

	StoreOrder(order *model.Order) (err error)
	GetOrderHistory(form GetOrderHistoryForm) (orders []model.Order, err error)
}
