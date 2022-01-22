package repository

import "fmt"

var (
	ErrNbElementTooLarge    = fmt.Errorf("cannot get enough last prices, number of element is too large")
	ErrNbValuesCannotBeZero = fmt.Errorf("nbValues must be greater than 0")
	ErrPairNotFound         = fmt.Errorf("pair cannot be found")
	ErrPlatformNotFound     = fmt.Errorf("platform prices cannot be found")
	ErrPriceNil             = fmt.Errorf("price is nil")
	ErrEntityNotFound       = fmt.Errorf("entity cannot be found in repository")
)
