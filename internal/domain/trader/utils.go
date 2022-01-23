package trader

import "crypto-bot/internal/constant/tradingConst"

var (
	currencyNeededByPairToBuy = map[string]string{
		tradingConst.BtcEurPair:  tradingConst.EuroCurrency,
		tradingConst.DashEurPair: tradingConst.EuroCurrency,
		tradingConst.EthEurPair:  tradingConst.EuroCurrency,
	}
	currencyNeededByPairToSell = map[string]string{
		tradingConst.BtcEurPair:  tradingConst.BtcCurrency,
		tradingConst.DashEurPair: tradingConst.DashCurrency,
		tradingConst.EthEurPair:  tradingConst.EthCurrency,
	}
)

// CurrencyNeededToTradePair will return currency that must be used to trade specific pair and specific tradeType (buy or sell)
func CurrencyNeededToTradePair(pair string, tradeType string) (currency string, err error) {
	switch tradeType {
	case tradingConst.BuySideOrder:
		currency, ok := currencyNeededByPairToBuy[pair]
		if !ok {
			return currency, errCannotFindCurrencyForPair(pair)
		}
	case tradingConst.SellSideOrder:
		currency, ok := currencyNeededByPairToSell[pair]
		if !ok {
			return currency, errCannotFindCurrencyForPair(pair)
		}
	default:
		return currency, errCannotFindCurrencyForPair(pair)
	}
	return
}
