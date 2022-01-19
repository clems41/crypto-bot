package repository

import "crypto-bot/internal/repository/repositoryModel"

type Price interface {
	GetLast(nbElement int, pair string) (prices []*repositoryModel.Price, err error)
	Store(price *repositoryModel.Price) (err error)
}

type Position interface {
	Store(position *repositoryModel.Position) (err error)
	Get(positionID string) (position *repositoryModel.Position, err error)
}
