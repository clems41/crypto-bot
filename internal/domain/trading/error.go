package trading

import "fmt"

var (
	errCurrencyNotInBalance = fmt.Errorf("pair is not define in this wallet")
	errPairNotFound         = fmt.Errorf("pair cannot be found")
)
