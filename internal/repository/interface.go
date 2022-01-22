package repository

import "crypto-bot/internal/repository/repositoryModel"

type Price interface {
	// GetMinimumAskPriceForNValues will return minimum price.Ask
	GetMinimumAskPriceForNValues(platformName string, pair string, nbValues int) (minimum float64, err error)
	GetCurrentPrice(platformName string, pair string) (price *repositoryModel.Price, err error)
	Store(price *repositoryModel.Price) (err error)
}

type Position interface {
	Store(position *repositoryModel.Position) (err error)
	Get(positionID string) (position *repositoryModel.Position, err error)
	GetOpenedPositions(platformName string, pairs ...string) (positions []*repositoryModel.Position, err error)
}

type Balance interface {
	Update(balance *repositoryModel.Balance) (err error)
	GetCurrentBalance(platformName string, currency string) (balance *repositoryModel.Balance, err error)
}
