package trading

import "fmt"

var (
	errCurrencyNotInBalance = fmt.Errorf("currency is not in wallet")
	errPairNotFound         = fmt.Errorf("pair cannot be found")
)
