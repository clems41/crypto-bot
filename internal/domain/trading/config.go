package trading

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envDelayBetweenEachRunInMilliSeconds     = "DELAY_RUN_MS"
	envMinimumResultInPercentToClosePosition = "MIN_RESULT"
	envNumberOfPreviousPricesToCompare       = "NB_PRICES"
	envMaxOpenedPositionsByPair              = "MAX_POSITION_PAIR"
	envMinimumAmountToOpenPosition           = "MIN_AMOUNT"
)

const (
	defaultDelayBetweenEachRunInMilliSeconds     = "5000"
	defaultMinimumResultInPercentToClosePosition = "-1"
	defaultNumberOfPreviousPricesToCompare       = "1"
	defaultMaxOpenedPositionsByPair              = "1"
	defaultMinimumAmountToOpenPosition           = "10"
)

var (
	pairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.DashEurPair,
		},
	}
)

func GetConfigFromEnvOrDefault() (config Config, err error) {
	delayBetweenEachRunInMilliSecondsStr := envUtils.GetFromEnvOrDefault(envDelayBetweenEachRunInMilliSeconds, defaultDelayBetweenEachRunInMilliSeconds)
	minimumResultInPercentToClosePositionStr := envUtils.GetFromEnvOrDefault(envMinimumResultInPercentToClosePosition, defaultMinimumResultInPercentToClosePosition)
	numberOfPreviousPricesToCompareStr := envUtils.GetFromEnvOrDefault(envNumberOfPreviousPricesToCompare, defaultNumberOfPreviousPricesToCompare)
	maxOpenedPositionsByPairStr := envUtils.GetFromEnvOrDefault(envMaxOpenedPositionsByPair, defaultMaxOpenedPositionsByPair)
	minimumAmountToOpenPositionStr := envUtils.GetFromEnvOrDefault(envMinimumAmountToOpenPosition, defaultMinimumAmountToOpenPosition)
	delayBetweenEachRunInMilliSeconds, err := strconv.Atoi(delayBetweenEachRunInMilliSecondsStr)
	if err != nil {
		return
	}
	minimumResultInPercentToClosePosition, err := strconv.ParseFloat(minimumResultInPercentToClosePositionStr, 64)
	if err != nil {
		return
	}
	numberOfPreviousPricesToCompare, err := strconv.Atoi(numberOfPreviousPricesToCompareStr)
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
		DelayBetweenEachRunInMilliSeconds:     delayBetweenEachRunInMilliSeconds,
		MinimumResultInPercentToClosePosition: minimumResultInPercentToClosePosition,
		NumberOfPreviousPricesToCompare:       numberOfPreviousPricesToCompare,
		MaxOpenedPositionsByPair:              maxOpenedPositionsByPair,
		MinimumAmountToOpenPosition:           minimumAmountToOpenPosition,
		PairToTradeByPlatform:                 pairsToTradeByPlatform,
	}
	return
}
