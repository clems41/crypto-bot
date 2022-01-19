package repositoryModel

import "time"

type Price struct {
	Date         time.Time
	PlatformName string
	Pair         string
	AskPrice     float64
	BidPrice     float64
}
