package tradingPlatform

type WalletView struct {
	Balance float32
}

type OpenPositionForm struct {
	Currency string
	Amount   float32
}

type OpenPositionView struct {
	PositionID string
}

type ClosePositionForm struct {
	PositionID string
}

type ClosePositionView struct {
	Result float32
}

type GetPriceForm struct {
	Currency string
}

type GetPriceView struct {
	Value float32
}
