package krakenApi

import (
	"crypto-bot/internal/constant/tradingConst"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

const (
	envKrakenApiSecret = "KRAKEN_API_SECRET"
	envKrakenApiKey    = "KRAKEN_API_KEY"
)

const (
	initBalance       = float64(100.0)
	fakeFeesInPercent = 0.17
)

var (
	pairConversion = map[string]string{
		tradingConst.BtcEurPair:  krakenapi.XXBTZEUR,
		tradingConst.DashEurPair: krakenapi.DASHEUR,
	}
)
