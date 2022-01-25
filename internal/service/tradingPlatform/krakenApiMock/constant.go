package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
)

const (
	envKrakenApiSecret = "KRAKEN_API_SECRET"
	envKrakenApiKey    = "KRAKEN_API_KEY"
)

const (
	takerFees = 0.26
	makerFees = 0.16
)

const (
	priceParameter          = "price"
	validateParameter       = "validate"
	closeOrderTypeParameter = "close[ordertype]"
	closePriceParameter     = "close[price]"
)

var (
	initialBalance = map[string]float64{
		tradingConst.EthCurrency:  0.04728132,
		tradingConst.EuroCurrency: 100,
		tradingConst.BtcCurrency:  0.00321543,
	}
)
