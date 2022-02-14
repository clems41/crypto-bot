package tradingUtils

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/mathUtils"
	"fmt"
)

// RemovePriceDecimal will return price with only usable decimal based on pair
func RemovePriceDecimal(price float64, pair tradingConst.Pair) (result float64, err error) {
	nbDecimal, ok := tradingConst.MaxDecimalByPair[pair]
	if !ok {
		return result, fmt.Errorf("cannot find decimal for pair %s", pair)
	}
	return mathUtils.RoundingAfterNDecimal(price, nbDecimal)
}

// RemoveAmountDecimalWithFloorRounding will return amount with only usable decimal based on currency.
// To avoid getting error 'insuffisant funds', amount should be floor rounded.
func RemoveAmountDecimalWithFloorRounding(amount float64, currency tradingConst.Currency) (result float64, err error) {
	nbDecimal, ok := tradingConst.MaxDecimalByCurrency[currency]
	if !ok {
		return result, fmt.Errorf("cannot find decimal for currency %s", currency)
	}
	return mathUtils.FloorRoundingAfterNDecimal(amount, nbDecimal)
}
