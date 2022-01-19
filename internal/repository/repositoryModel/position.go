package repositoryModel

import "time"

type Position struct {
	ID           string
	PlatformName string
	AskDate      time.Time
	BidDate      time.Time
	Pair         string
	Amount       float64
	AskPrice     float64
	BidPrice     float64
	Result       float64
	Closed       bool
}
