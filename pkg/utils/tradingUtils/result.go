package tradingUtils

import (
	"math"
	"time"
)

// EstimateProfit will estimate money that will be earned for a specific duration based on actual gain on a defined period of time
func EstimateProfit(initialBalance, finalBalance float64, tradingDuration, durationToEstimate time.Duration) (
	profit float64) {
	if initialBalance != 0 {
		gainCoefficient := finalBalance / initialBalance
		//gainCoefficient := 1 + (finalBalance-initialBalance)/initialBalance
		if tradingDuration != 0 {
			durationRatio := durationToEstimate.Seconds() / tradingDuration.Seconds()
			profit = initialBalance * math.Pow(gainCoefficient, durationRatio)
		}
	}
	return
}

// GetProfit will return value earned with this transaction with initial amount, it only can be positive because you can't sell less than 0
func GetProfit(askPrice, bidPrice, amount float64) (profit float64) {
	profit = bidPrice * amount / askPrice
	return
}

// GetResult will return value earned with this transaction without initial amount, can be positive or negative
func GetResult(askPrice, bidPrice, amount float64) (result float64) {
	profit := GetProfit(askPrice, bidPrice, amount)
	result = profit - amount
	return
}

// GetResultInPercent will return value earned in percent based on initial amount
func GetResultInPercent(askPrice, bidPrice, amount float64) (percent float64) {
	result := GetResult(askPrice, bidPrice, amount)
	if result != 0 {
		percent = result / amount * 100
	}
	return
}
