package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
)

var _ repository.Price = (*priceRepo)(nil)

type priceRepo struct {
	prices map[string][]*repositoryModel.Price // store all prices by currency
}

func NewPriceRepository() (*priceRepo, error) {
	return &priceRepo{
		prices: make(map[string][]*repositoryModel.Price),
	}, nil
}

func (repo *priceRepo) GetLast(nbElement int, currency string) (prices []*repositoryModel.Price, err error) {
	prices, ok := repo.prices[currency]
	if !ok || len(prices) == 0 {
		return nil, errCurrencyNotFound
	}
	if len(prices) < nbElement {
		return nil, errNbElementTooLarge
	}
	return prices[len(prices)-nbElement:], nil
}

func (repo *priceRepo) Store(price *repositoryModel.Price) (err error) {
	if price == nil {
		return errPriceNil
	}
	repo.prices[price.Currency] = append(repo.prices[price.Currency], price)
	return
}
