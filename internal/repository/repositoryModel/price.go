package repositoryModel

import "time"

type Price struct {
	Date     time.Time
	Currency string
	AskPrice float64
	BidPrice float64
}
