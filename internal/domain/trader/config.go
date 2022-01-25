package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envDelayBetweenEachRunInMilliSeconds      = "DELAY_RUN_MS"
	envMinimumResultInPercentBeforeCloseOrder = "MIN_RESULT"
)

const (
	defaultDelayBetweenEachRunInMilliSeconds      = "5000"
	defaultMinimumResultInPercentBeforeCloseOrder = "0.7"
)

var (
	pairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenMockPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.EthEurPair,
			tradingConst.EthBtcPair,
		},
	}
)

func GetConfigFromEnvOrDefault() (config Config, err error) {
	delayBetweenEachRunInMilliSecondsStr := envUtils.GetFromEnvOrDefault(envDelayBetweenEachRunInMilliSeconds, defaultDelayBetweenEachRunInMilliSeconds)
	minimumResultInPercentBeforeCloseOrderStr := envUtils.GetFromEnvOrDefault(envMinimumResultInPercentBeforeCloseOrder, defaultMinimumResultInPercentBeforeCloseOrder)
	delayBetweenEachRunInMilliSeconds, err := strconv.Atoi(delayBetweenEachRunInMilliSecondsStr)
	if err != nil {
		return
	}
	minimumResultInPercentBeforeCloseOrder, err := strconv.ParseFloat(minimumResultInPercentBeforeCloseOrderStr, 64)
	if err != nil {
		return
	}
	config = Config{
		DelayBetweenEachRunInMilliSeconds:     delayBetweenEachRunInMilliSeconds,
		MinimumResultInPercentToClosePosition: minimumResultInPercentBeforeCloseOrder,
		PairToTradeByPlatform:                 pairsToTradeByPlatform,
	}
	return
}
