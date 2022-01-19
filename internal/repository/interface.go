package repository

import "crypto-bot/internal/repository/repositoryModel"

type Price interface {
	GetLast(nbElement int, pair string) (prices []*repositoryModel.Price, err error)
	Store(price *repositoryModel.Price) (err error)
}
