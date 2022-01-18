package tradingPlatform

import "time"

/* Form */

type OpenPositionForm struct {
	Currency string
	Amount   float64
	AskPrice float64
}

type ClosePositionForm struct {
	PositionID string
}

type GetPriceForm struct {
	Currency string
}

/* View */

type PositionView struct {
	ID       string
	Currency string
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
	Result float64
}

type GetPriceView struct {
	AskPrice float64
	BidPrice float64
	Date     time.Time
}
