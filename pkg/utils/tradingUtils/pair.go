package tradingUtils

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/mathUtils"
	"fmt"
)

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
		tradingConst.EthBtcPair: {
			tradingConst.BuySideOrder:  tradingConst.BtcCurrency,
			tradingConst.SellSideOrder: tradingConst.EthCurrency,
		},
	}
)

var (
	maxDecimalByPair = map[string]int{
		tradingConst.BtcEurPair:  1,
		tradingConst.DashEurPair: 1,
		tradingConst.EthEurPair:  1,
		tradingConst.EthBtcPair:  1,
	}
)

// CurrencyNeededToTradePair will return currency that must be used to trade specific pair and specific orderSide (buy or sell).
// Ex : for ETH/EUR and buy --> EUR
// Ex : for BTC/EUR and sell --> BTC
func CurrencyNeededToTradePair(pair string, orderSide string) (currency string, ok bool) {
	currency, ok = currencyNeededByPairBySide[pair][orderSide]
	return
}

// CurrencyGotAfterTradingPair will return currency that will be traded with this pair and order side (buy or sell).
// Ex : for ETH/EUR and buy --> ETH
// Ex : for BTC/EUR and sell --> EUR
func CurrencyGotAfterTradingPair(pair string, orderSide string) (currency string, ok bool) {
	if orderSide == tradingConst.BuySideOrder {
		orderSide = tradingConst.SellSideOrder
	} else {
		orderSide = tradingConst.BuySideOrder
	}
	currency, ok = currencyNeededByPairBySide[pair][orderSide]
	return
}

func RemovePriceDecimal(price float64, pair string) (result float64, err error) {
	nbDecimal, ok := maxDecimalByPair[pair]
	if !ok {
		return result, fmt.Errorf("cannot find decimal for pair %s", pair)
	}
	return mathUtils.RemoveNDecimal(price, nbDecimal)
}
