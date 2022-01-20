package trading

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/repository/localRepository"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/internal/service/tradingPlatform/tradingPlatformMock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAmountToInvest(t *testing.T) {
	tradingMock, err := tradingPlatformMock.New()
	require.NoError(t, err)
	priceRepo, err := localRepository.NewPriceRepository()
	require.NoError(t, err)
	positionRepo, err := localRepository.NewPositionRepository()
	require.NoError(t, err)
	pairToTrade := map[string][]string{
		tradingConst.MockPlatform: {
			tradingConst.DashEurPair,
			tradingConst.BtcEurPair,
			tradingConst.BtcUsdPair,
		},
	}
	config := Config{
		MaxOpenedPositionsByPair:    2,
		MinimumAmountToOpenPosition: 10,
		PairToTradeByPlatform:       pairToTrade,
	}
	_, err = NewService(
		&config,
		[]tradingPlatform.Api{tradingMock},
		priceRepo,
		positionRepo,
	)
	require.NoError(t, err)

	// Store opened positions

}
