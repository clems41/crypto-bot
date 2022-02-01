package krakenApi

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
	assetPairConverter = map[string]string{
		tradingConst.BtcEurPair:  "XBTEUR",
		tradingConst.DashEurPair: "DASHEUR",
		tradingConst.EthEurPair:  "ETHEUR",
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

var (
	closeTypeDescriptionConverter = map[string]string{
		tradingConst.LimitOrderType:      "limit",
		tradingConst.StopLossOrderType:   "stop loss",
		tradingConst.TakeProfitOrderType: "take profit",
	}
)

var (
	statusConverter = map[string]string{
		tradingConst.OpenOrderStatus:   "open",
		tradingConst.CloseOrderStatus:  "closed",
		tradingConst.CancelOrderStatus: "canceled",
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

func GetProjectAssetPair(krakenPair string) (projectPair string, err error) {
	for project, kraken := range assetPairConverter {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, fmt.Errorf("cannot find project pair for %s", krakenPair)
}

func GetProjectSide(krakenSide string) (projectSide string, err error) {
	for project, kraken := range sideConverter {
		if kraken == krakenSide {
			projectSide = project
			return
		}
	}
	return projectSide, fmt.Errorf("cannot find side for %s", krakenSide)
}

func GetProjectOrderType(krakenType string) (projectType string, err error) {
	for project, kraken := range typeConverter {
		if kraken == krakenType {
			projectType = project
			return
		}
	}
	return projectType, fmt.Errorf("cannot find type for %s", krakenType)
}

func GetProjectStatus(krakenStatus string) (projectStatus string, err error) {
	for project, kraken := range statusConverter {
		if kraken == krakenStatus {
			projectStatus = project
			return
		}
	}
	return projectStatus, fmt.Errorf("cannot find status for %s", krakenStatus)
}
