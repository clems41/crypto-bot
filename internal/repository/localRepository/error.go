package localRepository

import "fmt"

var (
	errNbElementTooLarge = fmt.Errorf("cannot get enough last prices, number of element is too large")
	errCurrencyNotFound  = fmt.Errorf("currency cannot be found")
	errPriceNil          = fmt.Errorf("price is nil")
)
