package repository

import "crypto-bot/internal/repository/repositoryModel"

type Price interface {
	GetLast(platformName string, nbElement int, pair string) (prices []*repositoryModel.Price, err error)
	GetCurrentPrice(platformName string, pair string) (price *repositoryModel.Price, err error)
	Store(price *repositoryModel.Price) (err error)
}

type Position interface {
	Store(position *repositoryModel.Position) (err error)
	Get(positionID string) (position *repositoryModel.Position, err error)
	GetOpenedPositions(platformName string, pairs ...string) (positions []*repositoryModel.Position, err error)
}
