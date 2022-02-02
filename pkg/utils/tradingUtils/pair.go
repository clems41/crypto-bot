package tradingUtils

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/mathUtils"
	"fmt"
)

var (
	currencyNeededByPairBySide = map[tradingConst.Pair]map[tradingConst.OrderSide]tradingConst.Currency{
		tradingConst.BitcoinEuro: {
			tradingConst.Buy:  tradingConst.Euro,
			tradingConst.Sell: tradingConst.Bitcoin,
		},
		tradingConst.DashEuro: {
			tradingConst.Buy:  tradingConst.Euro,
			tradingConst.Sell: tradingConst.Dash,
		},
		tradingConst.EthereumEuro: {
			tradingConst.Buy:  tradingConst.Euro,
			tradingConst.Sell: tradingConst.Ethereum,
		},
		tradingConst.EthereumBitcoin: {
			tradingConst.Buy:  tradingConst.Bitcoin,
			tradingConst.Sell: tradingConst.Ethereum,
		},
		tradingConst.CardanoEuro: {
			tradingConst.Buy:  tradingConst.Euro,
			tradingConst.Sell: tradingConst.Cardano,
		},
	}
)

var (
	maxDecimalByPair = map[tradingConst.Pair]int{
		tradingConst.BitcoinEuro:     1,
		tradingConst.DashEuro:        1,
		tradingConst.EthereumEuro:    1,
		tradingConst.EthereumBitcoin: 1,
		tradingConst.CardanoEuro:     5,
	}
)

// CurrencyNeededToTradePair will return currency that must be used to trade specific pair and specific orderSide (buy or sell).
// Ex : for ETH/EUR and buy --> EUR
// Ex : for BTC/EUR and sell --> BTC
func CurrencyNeededToTradePair(pair tradingConst.Pair, orderSide tradingConst.OrderSide) (currency tradingConst.Currency, ok bool) {
	currency, ok = currencyNeededByPairBySide[pair][orderSide]
	return
}

// CurrencyGotAfterTradingPair will return currency that will be traded with this pair and order side (buy or sell).
// Ex : for ETH/EUR and buy --> ETH
// Ex : for BTC/EUR and sell --> EUR
func CurrencyGotAfterTradingPair(pair tradingConst.Pair, orderSide tradingConst.OrderSide) (currency tradingConst.Currency, ok bool) {
	if orderSide == tradingConst.Buy {
		orderSide = tradingConst.Sell
	} else {
		orderSide = tradingConst.Buy
	}
	currency, ok = currencyNeededByPairBySide[pair][orderSide]
	return
}

func RemovePriceDecimal(price float64, pair tradingConst.Pair) (result float64, err error) {
	nbDecimal, ok := maxDecimalByPair[pair]
	if !ok {
		return result, fmt.Errorf("cannot find decimal for pair %s", pair)
	}
	return mathUtils.RemoveNDecimal(price, nbDecimal)
}
