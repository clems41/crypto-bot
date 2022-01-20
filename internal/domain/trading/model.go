package trading

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
	// NumberOfPreviousPricesToCompare defines number of previous prices to use with opening position algorithm
	NumberOfPreviousPricesToCompare int
	// MaxOpenedPositionsByPair defines max number of position that should be opened for specific pair
	MaxOpenedPositionsByPair int
	// MinimumAmountToOpenPosition defines minimum value that should be used to open new position, if less than this value position will not be open
	MinimumAmountToOpenPosition float64
	// PairToTradeByPlatform defines all pairs that should trade for each platform
	PairToTradeByPlatform map[string][]string
}
