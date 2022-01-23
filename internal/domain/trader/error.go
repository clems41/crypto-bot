package trader

import "fmt"

func errCurrencyNotInBalance(currency string) error {
	return fmt.Errorf("currency %s is not in wallet", currency)
}

func errPlatformNotExist(platformName string) error {
	return fmt.Errorf("platform %s doesn't exist", platformName)
}

func errCannotFindCurrencyForPair(pair string) error {
	return fmt.Errorf("currency needed to trade pair %s cannot be found ", pair)
}
