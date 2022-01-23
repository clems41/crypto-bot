package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

const (
	envKrakenApiSecret = "KRAKEN_API_SECRET"
	envKrakenApiKey    = "KRAKEN_API_KEY"
)

const (
	initBalance = float64(100.0)
	takerFees   = 0.26
	makerFees   = 0.16
)

var (
	currencyConverter = map[string]string{
		"ZEUR": tradingConst.EuroCurrency,
		"DASH": tradingConst.DashCurrency,
		"XXBT": tradingConst.BtcCurrency,
		"XETH": tradingConst.EthCurrency,
	}
)

var (
	pairConverter = map[string]string{
		tradingConst.BtcEurPair:  krakenapi.XXBTZEUR,
		tradingConst.DashEurPair: krakenapi.DASHEUR,
		tradingConst.BtcUsdPair:  krakenapi.XXBTZUSD,
		tradingConst.EthEurPair:  krakenapi.XETHZEUR,
	}
)
