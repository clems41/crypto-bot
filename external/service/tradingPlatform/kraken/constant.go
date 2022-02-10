package kraken

import "time"

const (
	EnvKrakenApiSecret = "KRAKEN_API_SECRET"
	EnvKrakenApiKey    = "KRAKEN_API_KEY"
)

const (
	TakerFees = 0.26
	MakerFees = 0.16
)

const (
	ValidateParameter        = "validate"
	PriceParameter           = "price"
	CloseOrderTypeParameter  = "close_order_type"
	ClosePriceParameter      = "close_price"
	StartCloseOrderParameter = "start"
)

const (
	NbRequestRetries    = 3
	DelayBetweenRetries = 10 * time.Second
)
