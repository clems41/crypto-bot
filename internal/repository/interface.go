package repository

type Repository interface {
	StorePrice(price *model.Price) (err error)
	StoreBalance(balance *model.Balance) (err error)
	StoreOrder(order *model.Order) (err error)
}
