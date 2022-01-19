package repositoryModel

import "time"

type Price struct {
	Date     time.Time
	Pair     string
	AskPrice float64
	BidPrice float64
}
