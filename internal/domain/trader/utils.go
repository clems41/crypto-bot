package trader

import "crypto-bot/internal/constant/tradingConst"

var (
	currencyNeededByPairBySide = map[string]map[string]string{
		tradingConst.BtcEurPair: {
			tradingConst.BuySideOrder:  tradingConst.EuroCurrency,
			tradingConst.SellSideOrder: tradingConst.BtcCurrency,
		},
		tradingConst.DashEurPair: {
			tradingConst.BuySideOrder:  tradingConst.EuroCurrency,
			tradingConst.SellSideOrder: tradingConst.DashCurrency,
		},
		tradingConst.EthEurPair: {
			tradingConst.BuySideOrder:  tradingConst.EuroCurrency,
			tradingConst.SellSideOrder: tradingConst.EthCurrency,
		},
	}
)

// CurrencyNeededToTradePair will return currency that must be used to trade specific pair and specific orderSide (buy or sell)
func CurrencyNeededToTradePair(pair string, orderSide string) (currency string, err error) {
	currency, ok := currencyNeededByPairBySide[pair][orderSide]
	if !ok {
		return currency, errCannotFindCurrencyForPair
	}
	return
}
