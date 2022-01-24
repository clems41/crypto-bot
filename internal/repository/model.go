package repository

import "time"

/* Get Price */

type GetPriceHistoryForm struct {
	PlatformName string
	Pair         string
	SinceTime    time.Time
}
