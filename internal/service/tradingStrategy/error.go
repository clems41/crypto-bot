package tradingStrategy

import "fmt"

var (
	errNotEnoughPrice = fmt.Errorf("cannot determine if order should be open, need more prices")
)
