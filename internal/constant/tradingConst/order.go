package tradingConst

type OrderSide string

const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

type OrderType string

const (
	None       OrderType = "none"
	Market     OrderType = "market"
	Limit      OrderType = "limit"
	StopLoss   OrderType = "stop-loss"
	TakeProfit OrderType = "take-profit"
)

type OrderStatus string

const (
	Open   OrderStatus = "open"
	Close  OrderStatus = "close"
	Cancel OrderStatus = "cancel"
)
