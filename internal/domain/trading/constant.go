package trading

import "crypto-bot/internal/constant/currencyConst"

const (
	nbPreviousValues                            = 10
	delayBeforeNewAlgoApplicationInMilliseconds = 2000
	gainPercentToClosePosition                  = 1.005 // meaning 1%
	nbOfIncreasingValueToOpenPosition           = 6
	maxOpenedPositions                          = 5
)

var (
	currenciesToTrade = []string{
		currencyConst.BtcEurPair,
	}
)
