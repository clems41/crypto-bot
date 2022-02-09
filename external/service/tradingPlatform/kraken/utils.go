package kraken

import (
	"crypto-bot/internal/constant/tradingConst"
	"fmt"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

var (
	CurrencyConverterFromKrakenToProject = map[string]tradingConst.Currency{
		"ZEUR": tradingConst.Euro,
		"DASH": tradingConst.Dash,
		"XXBT": tradingConst.Bitcoin,
		"XETH": tradingConst.Ethereum,
		"ADA":  tradingConst.Cardano,
	}
)

var (
	PairConverter = map[tradingConst.Pair]string{
		tradingConst.BitcoinEuro:     krakenapi.XXBTZEUR,
		tradingConst.DashEuro:        krakenapi.DASHEUR,
		tradingConst.BitcoinUSDollar: krakenapi.XXBTZUSD,
		tradingConst.EthereumEuro:    krakenapi.XETHZEUR,
		tradingConst.EthereumBitcoin: krakenapi.XETHXXBT,
		tradingConst.CardanoEuro:     krakenapi.ADAEUR,
	}
)

var (
	AssetPairConverter = map[tradingConst.Pair]string{
		tradingConst.BitcoinEuro:  "XBTEUR",
		tradingConst.DashEuro:     "DASHEUR",
		tradingConst.EthereumEuro: "ETHEUR",
		tradingConst.CardanoEuro:  "ADAEUR",
	}
)

var (
	SideConverter = map[tradingConst.OrderSide]string{
		tradingConst.Buy:  "buy",
		tradingConst.Sell: "sell",
	}
)

var (
	TypeConverter = map[tradingConst.OrderType]string{
		tradingConst.Market:     krakenapi.OTMarket,
		tradingConst.Limit:      krakenapi.OTLimit,
		tradingConst.StopLoss:   krakenapi.OTStopLoss,
		tradingConst.TakeProfit: krakenapi.OTTakeProfi,
		tradingConst.None:       "",
	}
)

var (
	CloseTypeDescriptionConverter = map[tradingConst.OrderType]string{
		tradingConst.Limit:      "limit",
		tradingConst.StopLoss:   "stop loss",
		tradingConst.TakeProfit: "take profit",
	}
)

var (
	StatusConverter = map[tradingConst.OrderStatus]string{
		tradingConst.Open:   "open",
		tradingConst.Close:  "closed",
		tradingConst.Cancel: "canceled",
	}
)

func GetKrakenPair(projectPair tradingConst.Pair) (krakenPair string, err error) {
	krakenPair, ok := PairConverter[projectPair]
	if !ok {
		return krakenPair, fmt.Errorf("cannot find kraken pair for %s", projectPair)
	}
	return
}

func GetProjectPair(krakenPair string) (projectPair tradingConst.Pair, err error) {
	for project, kraken := range PairConverter {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, fmt.Errorf("cannot find project pair for %s", krakenPair)
}

func GetProjectAssetPair(krakenPair string) (projectPair tradingConst.Pair, err error) {
	for project, kraken := range AssetPairConverter {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, fmt.Errorf("cannot find project pair for %s", krakenPair)
}

func GetProjectSide(krakenSide string) (projectSide tradingConst.OrderSide, err error) {
	for project, kraken := range SideConverter {
		if kraken == krakenSide {
			projectSide = project
			return
		}
	}
	return projectSide, fmt.Errorf("cannot find side for %s", krakenSide)
}

func GetProjectOrderType(krakenType string) (projectType tradingConst.OrderType, err error) {
	for project, kraken := range TypeConverter {
		if kraken == krakenType {
			projectType = project
			return
		}
	}
	return projectType, fmt.Errorf("cannot find type for %s", krakenType)
}

func GetProjectStatus(krakenStatus string) (projectStatus tradingConst.OrderStatus, err error) {
	for project, kraken := range StatusConverter {
		if kraken == krakenStatus {
			projectStatus = project
			return
		}
	}
	return projectStatus, fmt.Errorf("cannot find status for %s", krakenStatus)
}
