package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"fmt"
	krakenapi "github.com/beldur/kraken-go-api-client"
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
		tradingConst.EthBtcPair:  krakenapi.XETHXXBT,
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

func GetKrakenPair(projectPair string) (krakenPair string, err error) {
	krakenPair, ok := pairConverter[projectPair]
	if !ok {
		return krakenPair, fmt.Errorf("cannot find kraken pair for %s", projectPair)
	}
	return
}

func GetProjectPair(krakenPair string) (projectPair string, err error) {
	for project, kraken := range pairConverter {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, fmt.Errorf("cannot find project pair for %s", krakenPair)
}
