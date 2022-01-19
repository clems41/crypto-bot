package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
)

var _ repository.Price = (*priceRepo)(nil)

type priceRepo struct {
	prices map[string]map[string][]*repositoryModel.Price // store all prices by platform and by pair
}

func NewPriceRepository() (*priceRepo, error) {
	return &priceRepo{
		prices: make(map[string]map[string][]*repositoryModel.Price),
	}, nil
}

func (repo *priceRepo) GetLast(platformName string, nbElement int, pair string) (result []*repositoryModel.Price, err error) {
	pricesByPair, ok := repo.prices[platformName]
	if !ok {
		return nil, repository.ErrPlatformNotFound
	}
	prices, ok := pricesByPair[pair]
	if !ok {
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
	if repo.prices[price.PlatformName] == nil {
		repo.prices[price.PlatformName] = make(map[string][]*repositoryModel.Price)
	}
	repo.prices[price.PlatformName][price.Pair] = append(repo.prices[price.PlatformName][price.Pair], price)
	return
}

func (repo *priceRepo) GetCurrentPrice(platformName string, pair string) (result *repositoryModel.Price, err error) {
	pricesByPair, ok := repo.prices[platformName]
	if !ok {
		return nil, repository.ErrPlatformNotFound
	}
	for pairPrice, prices := range pricesByPair {
		if pair == pairPrice {
			if len(prices) > 0 {
				result = prices[len(prices)-1]
				break
			}
		}
	}
	return
}
