package tradingPlatform

import "time"

/* Add Order */

/* Cancel Order */

/* Cancel All Orders */

type CancelAllOrdersView struct {
	Count int
}

/* Get Price */

type GetPricesForm struct {
	Pair      string
	SinceTime time.Time
}

type GetPriceView struct {
}

/* Get Open Orders */

type GetOpenOrdersView struct {
}

/* Get All Orders */

type GetAllOrdersView struct {
}
