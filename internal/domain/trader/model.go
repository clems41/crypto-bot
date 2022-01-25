package trader

type Position struct {
	ID           string
	PlatformName string
	Pair         string
	Amount       float64
	AskPrice     float64
}
