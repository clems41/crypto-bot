package tradingPlatform

import "time"

type PositionView struct {
	ID       string
	Currency string
	Amount   float64
	Price    float64
	Result   float64
}

type GetOpenedPositionsView struct {
	Positions []PositionView
}

type GetAllPositionsView struct {
	Positions []PositionView
}

type WalletView struct {
	BalanceByCurrency map[string]float64
}

type OpenPositionForm struct {
	Currency string
	Amount   float64
}

type OpenPositionView struct {
	PositionID string
}

type ClosePositionForm struct {
	PositionID string
}

type ClosePositionView struct {
	Result float64
}

type GetPriceForm struct {
	Currency string
}

type GetPriceView struct {
	AskPrice float64
	BidPrice float64
	Date     time.Time
}
