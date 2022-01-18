package trading

type Position struct {
	ID       string
	Amount   float64
	AskPrice float64
	BidPrice float64
}

type Price struct {
	Timestamp uint
	AskPrice  float64
	BidPrice  float64
}
