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
	closeOrderTypeParameter = "close_order_type"
	closePriceParameter     = "close_price"
)

var (
	initialBalance = map[string]float64{
		tradingConst.EuroCurrency: 100,
	}
)
