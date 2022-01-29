package repository

import "time"

/* Get Price */

type GetPriceHistoryForm struct {
	PlatformName string
	Pair         string
	SinceTime    time.Time
}

/* Get Order */

type GetOrderHistoryForm struct {
	PlatformName string
	Pair         string
	Status       string
}
