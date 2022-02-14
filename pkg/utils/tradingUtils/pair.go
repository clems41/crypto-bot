package tradingUtils

import (
	"crypto-bot/internal/constant/tradingConst"
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
