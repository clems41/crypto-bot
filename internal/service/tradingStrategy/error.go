package tradingStrategy

import "fmt"

var (
	ErrNotEnoughPrice = fmt.Errorf("cannot determine if order should be open, need more prices")
)
