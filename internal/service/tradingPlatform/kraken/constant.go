package kraken

import (
	"crypto-bot/internal/constant/currencyConst"
	krakenapi "github.com/beldur/kraken-go-api-client"
)

const (
	apiSecret = "Jk9pEPFgt91kPyNs0ZSYvoObcD6iGEzsUhbzCXksetc+hXjIMMnr3eyYMq0SK4/OOSmENeJHMMnHug9/FmaBEA=="
	apiKey    = "4TBOIwtlVtOLPWBGOzzcao2tOheSXDwRu3ihZ2o73Fdhap+6oYirTM9d"
)

const (
	initBalance = float64(100.0)
)

var (
	currencyConversion = map[string]string{
		currencyConst.BtcEurPair: krakenapi.XXBTZEUR,
	}
)
