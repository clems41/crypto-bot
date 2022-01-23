package maxOrMinAlgo

import (
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
	defaultMinimumResultInPercentToClosePosition = "0.5"
	defaultNumberOfPreviousPricesToCompare       = "720" // 720 * 5000ms --> 1 hour
	defaultMaxOpenedPositionsByPair              = "1"
	defaultMinimumAmountToOpenPosition           = "10"
)

type Config struct {
	// DelayBetweenEachRun defines delay in milliseconds to wait before running new algorithm iteration
	DelayBetweenEachRunInMilliSeconds int
	// MinimumResultInPercentToClosePosition defines minimum result in percent to close position, position will not be closed if actual result is less than this value
	MinimumResultInPercentToClosePosition float64
	// NumberOfPreviousPricesToCompare defines number of previous prices to use with opening position algorithm
	NumberOfPreviousPricesToCompare int
	// MaxOpenedPositionsByPair defines max number of position that should be opened for specific pair
	MaxOpenedPositionsByPair int
	// MinimumAmountToOpenPosition defines minimum value that should be used to open new position, if less than this value position will not be open
	MinimumAmountToOpenPosition float64
}

func GetConfigFromEnvOrDefault() (config *Config, err error) {
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
	config = &Config{
		DelayBetweenEachRunInMilliSeconds:     delayBetweenEachRunInMilliSeconds,
		MinimumResultInPercentToClosePosition: minimumResultInPercentToClosePosition,
		NumberOfPreviousPricesToCompare:       numberOfPreviousPricesToCompare,
		MaxOpenedPositionsByPair:              maxOpenedPositionsByPair,
		MinimumAmountToOpenPosition:           minimumAmountToOpenPosition,
	}
	return
}
