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
	takerFees = 0.26
	makerFees = 0.16
)

const (
	priceParameter    = "price"
	validateParameter = "validate"
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

var (
	sideConverter = map[string]string{
		tradingConst.BuySideOrder:  "buy",
		tradingConst.SellSideOrder: "sell",
	}
)

var (
	typeConverter = map[string]string{
		tradingConst.MarketOrderType:     krakenapi.OTMarket,
		tradingConst.LimitOrderType:      krakenapi.OTLimit,
		tradingConst.StopLossOrderType:   krakenapi.OTStopLoss,
		tradingConst.TakeProfitOrderType: krakenapi.OTTakeProfi,
	}
)
