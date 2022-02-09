package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
	"time"
)

const (
	envNumberOfPreviousPricesToCompare        = "NB_PRICES"
	envPercentPriceBelowToBuy                 = "PERCENT_BEFORE_BUY"
	envMinimumResultInPercentBeforeCloseOrder = "MIN_RESULT"
	envMaxOpenedOrdersByPair                  = "MAX_ORDER_PAIR"
	envMinimumAmount                          = "MIN_AMOUNT"
	envDelayBetweenEachRunInSeconds           = "DELAY_RUN_SEC"
)

const (
	defaultNumberOfPreviousPricesToCompare        = "180" // 180 * 10s --> 30min
	defaultPercentPriceBelowToBuy                 = "0.05"
	defaultMinimumResultInPercentBeforeCloseOrder = "0.1"
	defaultMaxOpenedOrdersByPair                  = "1"
	defaultMinimumAmount                          = "10"
	defaultDelayBetweenEachRunInSeconds           = "10"
)

var (
	initialPairsToTradeByPlatform = map[string][]tradingConst.Pair{
		tradingConst.KrakenMockPlatform: {
			tradingConst.BitcoinEuro,
			tradingConst.EthereumEuro,
			tradingConst.DashEuro,
			tradingConst.CardanoEuro,
		},
		tradingConst.KrakenPlatform: {
			tradingConst.BitcoinEuro,
			tradingConst.EthereumEuro,
			tradingConst.DashEuro,
			tradingConst.CardanoEuro,
		},
	}
)

type config struct {
	// NumberOfPreviousPricesToCompare defines number of previous prices to use with opening position algorithm
	NumberOfPreviousPricesToCompare int
	// PercentPriceBelowToBuy defines percent of price decrease needed to trigger order.
	// If index price = 100 and PercentPriceBelowToBuy = 50, order will be open with price = 100 * 50 / 100 = 50
	PercentPriceBelowToBuy float64
	// MinimumResultInPercentToClosePosition defines minimum result in percent to close position, position will not be closed if actual result is less than this value.
	// This value already includes platform fees. If equal to 0.1 and fees = 0.26, order will be close when result is at least : 1.0026 * 1.0026 * 1.001 = 1.0062 --> 0.62%
	MinimumResultInPercentToClosePosition float64
	// MaxOpenedOrdersByPair return limit of opened orders by pair
	MaxOpenedOrdersByPair int
	// MinimumAmount defines minimal amount to open order
	MinimumAmount float64
	// DelayBetweenEachRun defines duration to wait between each algorithm run
	DelayBetweenEachRun time.Duration
	// InitialPairsToTradeByPlatform defines all pairs that should be trade by platform
	InitialPairsToTradeByPlatform map[string][]tradingConst.Pair
}

func GetConfigFromEnvOrDefault() (cfg *config, err error) {
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
	delayBetweenEachRunInSecondsStr := envUtils.GetFromEnvOrDefault(envDelayBetweenEachRunInSeconds, defaultDelayBetweenEachRunInSeconds)
	delayBetweenEachRunInSeconds, err := strconv.Atoi(delayBetweenEachRunInSecondsStr)
	if err != nil {
		return
	}
	cfg = &config{
		NumberOfPreviousPricesToCompare:       numberOfPreviousPricesToCompare,
		PercentPriceBelowToBuy:                percentPriceBelowToBuy,
		MinimumResultInPercentToClosePosition: minimumResultInPercentBeforeCloseOrder,
		MaxOpenedOrdersByPair:                 maxOpenedOrdersByPair,
		MinimumAmount:                         minimumAmount,
		InitialPairsToTradeByPlatform:         initialPairsToTradeByPlatform,
		DelayBetweenEachRun:                   time.Duration(delayBetweenEachRunInSeconds) * time.Second,
	}
	logger.Infof("Following config will be used for trading algo : %+v", cfg)
	return
}
