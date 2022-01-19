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
	prices, err := repo.GetLast(platformName, 1, pair)
	require.Error(t, err)

	// Store some prices with correct pair, and should get them all
	var currencyPrices []repositoryModel.Price
	for range fakeData.FakeRange(5, 15) {
		price := repositoryModel.Price{
			PlatformName: platformName,
			Date:         time.Now(),
			Pair:         pair,
			AskPrice:     float64(fakeData.FakeIntBetween(1, 10000)),
			BidPrice:     float64(fakeData.FakeIntBetween(1, 10000)),
		}
		err = repo.Store(&price)
		require.NoError(t, err)
		currencyPrices = append(currencyPrices, price)
	}

	// Store some prices with incorrect pair, and should not get them
	var incorrectCurrencyPrices []repositoryModel.Price
	for range fakeData.FakeRange(5, 15) {
		price := repositoryModel.Price{
			Date:     time.Now(),
			Pair:     fake.Word(),
			AskPrice: float64(fakeData.FakeIntBetween(10000, 100000)),
			BidPrice: float64(fakeData.FakeIntBetween(10000, 100000)),
		}
		err = repo.Store(&price)
		require.NoError(t, err)
		incorrectCurrencyPrices = append(incorrectCurrencyPrices, price)
	}

	// Get all correct pair prices
	prices, err = repo.GetLast(platformName, len(currencyPrices), pair)
	require.NoError(t, err)
	require.Len(t, prices, len(currencyPrices))
	for _, price := range prices {
		require.Equal(t, pair, price.Pair)
		require.True(t, price.AskPrice > 1 && price.AskPrice < 10000,
			"Bitcoin ask price should be between 1 and 10 000 / unit but it is %0.2f", price.AskPrice)
		require.True(t, price.BidPrice > 1 && price.BidPrice < 10000,
			"Bitcoin bid price should be between 1 and 10 000 / unit but it is %0.2f", price.BidPrice)
	}
}
