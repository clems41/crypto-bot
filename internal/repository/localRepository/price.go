package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/pkg/utils/csvUtils"
	"crypto-bot/pkg/utils/pathUtils"
	"fmt"
	"time"
)

var _ repository.Price = (*priceRepo)(nil)

type priceRepo struct {
	prices      map[string]map[string][]*repositoryModel.Price // store all prices by platform and by pair
	csvFilePath string
}

func NewPriceRepository() (repo *priceRepo, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}
	csvPath := fmt.Sprintf("%s/%s", rootPath, csvFileNamePrice)

	// empty file or creating new one
	err = csvUtils.ForceCreateFile(csvPath)
	if err != nil {
		return
	}

	// adding column names
	columnNames := []string{
		"Date",
		"PlatformName",
		"Pair",
		"Ask",
		"Bid",
	}
	err = csvUtils.AppendLines(csvPath, columnNames)
	if err != nil {
		return
	}

	repo = &priceRepo{
		prices:      make(map[string]map[string][]*repositoryModel.Price),
		csvFilePath: csvPath,
	}
	return
}

func (repo *priceRepo) GetMinimumAskPriceForNValues(platformName string, pair string, nbValues int) (minimum float64, err error) {
	if nbValues <= 0 {
		return minimum, repository.ErrNbValuesCannotBeZero
	}
	// Get n prices for specific platform and pair
	pricesByPair, ok := repo.prices[platformName]
	if !ok {
		return minimum, repository.ErrPlatformNotFound
	}
	prices, ok := pricesByPair[pair]
	if !ok {
		return minimum, repository.ErrPairNotFound
	}
	if nbValues > len(prices) {
		return minimum, repository.ErrNbValuesTooLarge
	}

	// Find minimum ask price among n values
	minimum = prices[len(prices)-1].Ask
	for _, price := range prices[len(prices)-nbValues:] {
		if price.Ask < minimum {
			minimum = price.Ask
		}
	}
	return
}

func (repo *priceRepo) Store(price *repositoryModel.Price) (err error) {
	if price == nil {
		return repository.ErrPriceNil
	}
	if repo.prices[price.PlatformName] == nil {
		repo.prices[price.PlatformName] = make(map[string][]*repositoryModel.Price)
	}
	repo.prices[price.PlatformName][price.Pair] = append(repo.prices[price.PlatformName][price.Pair], price)

	// creating new csv line
	newLine := []string{
		price.Date.Format(time.RFC3339),
		price.PlatformName,
		price.Pair,
		fmt.Sprintf("%0.2f", price.Ask),
		fmt.Sprintf("%0.2f", price.Bid),
	}

	// update line or append it if not exists
	err = csvUtils.AppendLines(repo.csvFilePath, newLine)
	if err != nil {
		return
	}
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
