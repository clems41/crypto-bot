package tradingUtils

import "crypto-bot/internal/constant/tradingConst"

var (
	currencyNeededByPair = map[string]string{
		tradingConst.BtcEurPair:  tradingConst.EuroCurrency,
		tradingConst.DashEurPair: tradingConst.EuroCurrency,
	}
)

// CanTradePairUsingCurrency return true if currency can be used to trade specific pair.
func CanTradePairUsingCurrency(currency, pair string) (can bool) {
	currencyNeeded := CurrencyNeededToTradePair(pair)
	return currencyNeeded == currency
}

// CurrencyNeededToTradePair will return currency that must be used to trade specific pair
func CurrencyNeededToTradePair(pair string) (currency string) {
	currency, _ = currencyNeededByPair[pair]
	return
}
