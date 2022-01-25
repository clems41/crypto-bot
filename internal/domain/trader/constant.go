package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"time"
)

const (
	delayBetweenEachRun = 5 * time.Second
)

var (
	initialPairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenMockPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.EthEurPair,
		},
	}
)
