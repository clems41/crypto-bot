package tradingPlatform

import "time"

/* Form */

type OpenPositionForm struct {
	Pair     string
	Amount   float64
	AskPrice float64
}

type ClosePositionForm struct {
	PositionID string
}

type GetPriceForm struct {
	Pairs []string
}

/* View */

type PositionView struct {
	ID       string
	Pair     string
	Amount   float64
	AskPrice float64
	BidPrice float64
	Result   float64
}

type WalletView struct {
	BalanceByCurrency map[string]float64
}

type GetOpenedPositionsView struct {
	Positions []PositionView
}

type GetAllPositionsView struct {
	Positions []PositionView
}

type OpenPositionView struct {
	PositionID string
}

type ClosePositionView struct {
	AskPrice float64
	BidPrice float64
	Result   float64
	Profit   float64
}

type PriceView struct {
	AskPrice float64
	BidPrice float64
	Date     time.Time
}

type GetPriceView struct {
	PriceByPair map[string]PriceView // price by pair
}
