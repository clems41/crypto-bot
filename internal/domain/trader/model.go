package trader

type Position struct {
	ID           string
	PlatformName string
	Pair         string
	Amount       float64
	AskPrice     float64
}

type Config struct {
	// DelayBetweenEachRun defines delay in milliseconds to wait before running new algorithm iteration
	DelayBetweenEachRunInMilliSeconds int
	// MinimumResultInPercentToClosePosition defines minimum result in percent to close position, position will not be closed if actual result is less than this value
	MinimumResultInPercentToClosePosition float64
	// IntervalToComparePricesInMinutes defines duration in minutes to get prices history before making the choice to add new order
	IntervalToComparePricesInMinutes int
	// PairToTradeByPlatform defines all pairs that should trade for each platform
	PairToTradeByPlatform map[string][]string
}
