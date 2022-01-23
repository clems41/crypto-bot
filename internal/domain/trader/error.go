package trader

import "fmt"

var (
	errCurrencyNotInBalance      = fmt.Errorf("currency is not in wallet")
	errPlatformNotExist          = fmt.Errorf("platform doesn't exist")
	errCannotFindCurrencyForPair = fmt.Errorf("currency needed to trade pair cannot be found ")
)
