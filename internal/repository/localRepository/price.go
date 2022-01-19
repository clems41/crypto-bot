package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
)

var _ repository.Price = (*priceRepo)(nil)

type priceRepo struct {
	prices map[string][]*repositoryModel.Price // store all prices by pair
}

func NewPriceRepository() (*priceRepo, error) {
	return &priceRepo{
		prices: make(map[string][]*repositoryModel.Price),
	}, nil
}

func (repo *priceRepo) GetLast(nbElement int, pair string) (prices []*repositoryModel.Price, err error) {
	prices, ok := repo.prices[pair]
	if !ok || len(prices) == 0 {
		return nil, repository.ErrPairNotFound
	}
	if len(prices) < nbElement {
		return nil, repository.ErrNbElementTooLarge
	}
	return prices[len(prices)-nbElement:], nil
}

func (repo *priceRepo) Store(price *repositoryModel.Price) (err error) {
	if price == nil {
		return repository.ErrPriceNil
	}
	repo.prices[price.Pair] = append(repo.prices[price.Pair], price)
	return
}
