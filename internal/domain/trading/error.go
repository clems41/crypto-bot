package trading

import "fmt"

var (
	errCurrencyNotInBalance = fmt.Errorf("currency is not in wallet")
	errPairNotFound         = fmt.Errorf("pair cannot be found")
	errPriceNotFound        = fmt.Errorf("price cannot be found")
	errPlatformNotFound     = fmt.Errorf("platform cannot be found")
)
