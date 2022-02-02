package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"time"
)

const (
	delayBetweenEachRun = 10 * time.Second
)

var (
	initialPairsToTradeByPlatform = map[string][]tradingConst.Pair{
		tradingConst.KrakenMockPlatform: {
			tradingConst.BitcoinEuro,
			tradingConst.EthereumBitcoin,
			tradingConst.DashEuro,
			tradingConst.CardanoEuro,
		},
		tradingConst.KrakenPlatform: {
			tradingConst.BitcoinEuro,
			tradingConst.EthereumBitcoin,
			tradingConst.DashEuro,
			tradingConst.CardanoEuro,
		},
	}
)
