package tradingPlatform

import "fmt"

var (
	ErrPairNotFound                = fmt.Errorf("pair doesn't exists")
	ErrCannotFindPriceFromResponse = fmt.Errorf("price cannot be found in api response")
	ErrPositionNotFound            = fmt.Errorf("position not found")
	ErrEmptyResponse               = fmt.Errorf("response is nil from api")
)
