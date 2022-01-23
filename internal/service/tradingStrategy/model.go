package tradingStrategy

/* Config */

type Config struct {
	// NumberPreviousPricesNeeded Number of previous prices needed to know if order should be open
	NumberPreviousPricesNeeded int
}

/* Should Add Order */

type ShouldAddOrderForm struct {
}

type ShouldAddOrderView struct {
}
