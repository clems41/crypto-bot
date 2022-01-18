package kraken

import (
	"crypto-bot/internal/constant/currencyConst"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

const (
	envKrakenApiSecret = "KRAKEN_API_SECRET"
	envKrakenApiKey    = "KRAKEN_API_KEY"
)

const (
	initBalance = float64(100.0)
)

var (
	currencyConversion = map[string]string{
		currencyConst.BtcEurPair: krakenapi.XXBTZEUR,
	}
)
