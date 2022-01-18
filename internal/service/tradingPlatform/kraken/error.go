package kraken

import "fmt"

var (
	errCurrencyNotFound            = fmt.Errorf("currency doesn't exists")
	errCannotFindPriceFromResponse = fmt.Errorf("price cannot be found in api response")
	errPositionNotFound            = fmt.Errorf("position not found")
)
