package repository

import "fmt"

var (
	ErrNbElementTooLarge = fmt.Errorf("cannot get enough last prices, number of element is too large")
	ErrPairNotFound      = fmt.Errorf("pair cannot be found")
	ErrPriceNil          = fmt.Errorf("price is nil")
)
