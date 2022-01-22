package repositoryModel

import "time"

type Price struct {
	Date         time.Time
	PlatformName string
	Pair         string
	Ask          float64
	Bid          float64
}
