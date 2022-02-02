package repository

import (
	"crypto-bot/internal/constant/tradingConst"
	"time"
)

/* Get Price */

type GetPriceHistoryForm struct {
	PlatformName string
	Pair         tradingConst.Pair
	SinceTime    time.Time
}

/* Get Order */

type GetOrderHistoryForm struct {
	PlatformName string
	Pair         tradingConst.Pair
	Status       tradingConst.OrderStatus
}
