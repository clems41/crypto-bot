package trader

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envDelayBetweenEachRunInMilliSeconds = "DELAY_RUN_MS"
	envIntervalToComparePricesInMinutes  = "INTERVAL_PRICES"
	envMaxOpenedPositionsByPair          = "MAX_POSITION_PAIR"
	envMinimumAmountToOpenPosition       = "MIN_AMOUNT"
)

const (
	defaultDelayBetweenEachRunInMilliSeconds = "5000"
	defaultIntervalToComparePricesInMinutes  = "60" // --> 1 hour
	defaultMaxOpenedPositionsByPair          = "1"
	defaultMinimumAmountToOpenPosition       = "10"
)

var (
	pairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.DashEurPair,
			tradingConst.EthEurPair,
		},
	}
)

func GetConfigFromEnvOrDefault() (config Config, err error) {
	delayBetweenEachRunInMilliSecondsStr := envUtils.GetFromEnvOrDefault(envDelayBetweenEachRunInMilliSeconds, defaultDelayBetweenEachRunInMilliSeconds)
	intervalToComparePricesInMinutesStr := envUtils.GetFromEnvOrDefault(envIntervalToComparePricesInMinutes, defaultIntervalToComparePricesInMinutes)
	maxOpenedPositionsByPairStr := envUtils.GetFromEnvOrDefault(envMaxOpenedPositionsByPair, defaultMaxOpenedPositionsByPair)
	minimumAmountToOpenPositionStr := envUtils.GetFromEnvOrDefault(envMinimumAmountToOpenPosition, defaultMinimumAmountToOpenPosition)
	delayBetweenEachRunInMilliSeconds, err := strconv.Atoi(delayBetweenEachRunInMilliSecondsStr)
	if err != nil {
		return
	}
	intervalToComparePricesInMinutes, err := strconv.Atoi(intervalToComparePricesInMinutesStr)
	if err != nil {
		return
	}
	maxOpenedPositionsByPair, err := strconv.Atoi(maxOpenedPositionsByPairStr)
	if err != nil {
		return
	}
	minimumAmountToOpenPosition, err := strconv.ParseFloat(minimumAmountToOpenPositionStr, 64)
	if err != nil {
		return
	}
	config = Config{
		DelayBetweenEachRunInMilliSeconds: delayBetweenEachRunInMilliSeconds,
		IntervalToComparePricesInMinutes:  intervalToComparePricesInMinutes,
		MaxOpenedPositionsByPair:          maxOpenedPositionsByPair,
		MinimumAmountToOpenPosition:       minimumAmountToOpenPosition,
		PairToTradeByPlatform:             pairsToTradeByPlatform,
	}
	return
}
