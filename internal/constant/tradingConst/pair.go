package tradingConst

const (
	BtcEurPair  = "BTC_EUR"
	DashEurPair = "DASH_EUR"
)

var (
	currencyNeededByPair = map[string]string{
		BtcEurPair:  EuroCurrency,
		DashEurPair: EuroCurrency,
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
