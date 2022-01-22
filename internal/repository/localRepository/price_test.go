package localRepository

import (
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPriceRepo_GetAndStore(t *testing.T) {
	repo, err := NewPriceRepository()
	require.NoError(t, err)
	pair := fake.Word()
	platformName := fake.Word()

	// try to get last prices without storing any, should return error
	_, err = repo.GetMinimumAskPriceForNValues(platformName, pair, 1)
	require.Error(t, err)

	// Store some prices with correct pair, and should get them all
	var currencyPrices []repositoryModel.Price
	for range fakeData.FakeRange(5, 15) {
		price := repositoryModel.Price{
			PlatformName: platformName,
			Date:         time.Now(),
			Pair:         pair,
			Ask:          float64(fakeData.FakeIntBetween(1, 10000)),
			Bid:          float64(fakeData.FakeIntBetween(1, 10000)),
		}
		err = repo.Store(&price)
		require.NoError(t, err)
		currencyPrices = append(currencyPrices, price)
	}

	// Store some prices with incorrect pair, and should not get them
	var incorrectCurrencyPrices []repositoryModel.Price
	for range fakeData.FakeRange(5, 15) {
		price := repositoryModel.Price{
			Date: time.Now(),
			Pair: fake.Word(),
			Ask:  float64(fakeData.FakeIntBetween(10000, 100000)),
			Bid:  float64(fakeData.FakeIntBetween(10000, 100000)),
		}
		err = repo.Store(&price)
		require.NoError(t, err)
		incorrectCurrencyPrices = append(incorrectCurrencyPrices, price)
	}

	// Get all correct pair prices
	minimumAskPrice, err := repo.GetMinimumAskPriceForNValues(platformName, pair, len(currencyPrices))
	require.NoError(t, err)
	require.True(t, minimumAskPrice > 1 && minimumAskPrice < 10000,
		"Minimum ask price should be between 1 and 10 000 / unit but it is %0.2f", minimumAskPrice)
}
