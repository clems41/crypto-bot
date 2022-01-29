package minOrMaxAlgo

import (
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envNumberOfPreviousPricesToCompare        = "NB_PRICES"
	envPercentPriceBelowToBuy                 = "PERCENT_BEFORE_BUY"
	envMinimumResultInPercentBeforeCloseOrder = "MIN_RESULT"
	envMaxOpenedOrdersByPair                  = "MAX_ORDER_PAIR"
	envMinimumAmount                          = "MIN_AMOUNT"
)

const (
	defaultNumberOfPreviousPricesToCompare        = "1" // 200 * 5ms --> 13 min
	defaultPercentPriceBelowToBuy                 = "-0.05"
	defaultMinimumResultInPercentBeforeCloseOrder = "-0.6"
	defaultMaxOpenedOrdersByPair                  = "1"
	defaultMinimumAmount                          = "10"
)

type Config struct {
	// NumberOfPreviousPricesToCompare defines number of previous prices to use with opening position algorithm
	NumberOfPreviousPricesToCompare int
	// PercentPriceBelowToBuy defines percent of price decrease needed to trigger order.
	// If index price = 100 and PercentPriceBelowToBuy = 50, order will be open with price = 100 * 50 / 100 = 50
	PercentPriceBelowToBuy float64
	// MinimumResultInPercentToClosePosition defines minimum result in percent to close position, position will not be closed if actual result is less than this value
	MinimumResultInPercentToClosePosition float64
	// MaxOpenedOrdersByPair return limit of opened orders by pair
	MaxOpenedOrdersByPair int
	// MinimumAmount defines minimal amount to open order
	MinimumAmount float64
}

func GetConfigFromEnvOrDefault() (config *Config, err error) {
	maxOpenedOrdersByPairStr := envUtils.GetFromEnvOrDefault(envMaxOpenedOrdersByPair, defaultMaxOpenedOrdersByPair)
	maxOpenedOrdersByPair, err := strconv.Atoi(maxOpenedOrdersByPairStr)
	if err != nil {
		return
	}
	numberOfPreviousPricesToCompareStr := envUtils.GetFromEnvOrDefault(envNumberOfPreviousPricesToCompare, defaultNumberOfPreviousPricesToCompare)
	numberOfPreviousPricesToCompare, err := strconv.Atoi(numberOfPreviousPricesToCompareStr)
	if err != nil {
		return
	}
	percentPriceBelowToBuyStr := envUtils.GetFromEnvOrDefault(envPercentPriceBelowToBuy, defaultPercentPriceBelowToBuy)
	percentPriceBelowToBuy, err := strconv.ParseFloat(percentPriceBelowToBuyStr, 64)
	if err != nil {
		return
	}
	minimumResultInPercentBeforeCloseOrderStr := envUtils.GetFromEnvOrDefault(envMinimumResultInPercentBeforeCloseOrder, defaultMinimumResultInPercentBeforeCloseOrder)
	minimumResultInPercentBeforeCloseOrder, err := strconv.ParseFloat(minimumResultInPercentBeforeCloseOrderStr, 64)
	if err != nil {
		return
	}
	minimumAmountStr := envUtils.GetFromEnvOrDefault(envMinimumAmount, defaultMinimumAmount)
	minimumAmount, err := strconv.ParseFloat(minimumAmountStr, 64)
	if err != nil {
		return
	}
	config = &Config{
		NumberOfPreviousPricesToCompare:       numberOfPreviousPricesToCompare,
		PercentPriceBelowToBuy:                percentPriceBelowToBuy,
		MinimumResultInPercentToClosePosition: minimumResultInPercentBeforeCloseOrder,
		MaxOpenedOrdersByPair:                 maxOpenedOrdersByPair,
		MinimumAmount:                         minimumAmount,
	}
	logger.Infof("Following config will be used for trading algo : %+v", config)
	return
}
