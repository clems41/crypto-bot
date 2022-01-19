package trading

import "fmt"

var (
	errCurrencyNotInBalance = fmt.Errorf("currency is wallet")
	errPairNotFound         = fmt.Errorf("pair cannot be found")
)
