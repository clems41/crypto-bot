package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envDelayBetweenEachRunInMilliSeconds = "DELAY_RUN_MS"
	envIntervalToComparePricesInMinutes  = "INTERVAL_PRICES"
)

const (
	defaultDelayBetweenEachRunInMilliSeconds = "5000"
	defaultIntervalToComparePricesInMinutes  = "60" // --> 1 hour
)

var (
	pairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenMockPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.EthEurPair,
		},
	}
)

func GetConfigFromEnvOrDefault() (config Config, err error) {
	delayBetweenEachRunInMilliSecondsStr := envUtils.GetFromEnvOrDefault(envDelayBetweenEachRunInMilliSeconds, defaultDelayBetweenEachRunInMilliSeconds)
	intervalToComparePricesInMinutesStr := envUtils.GetFromEnvOrDefault(envIntervalToComparePricesInMinutes, defaultIntervalToComparePricesInMinutes)
	delayBetweenEachRunInMilliSeconds, err := strconv.Atoi(delayBetweenEachRunInMilliSecondsStr)
	if err != nil {
		return
	}
	intervalToComparePricesInMinutes, err := strconv.Atoi(intervalToComparePricesInMinutesStr)
	if err != nil {
		return
	}
	config = Config{
		DelayBetweenEachRunInMilliSeconds: delayBetweenEachRunInMilliSeconds,
		IntervalToComparePricesInMinutes:  intervalToComparePricesInMinutes,
		PairToTradeByPlatform:             pairsToTradeByPlatform,
	}
	return
}
