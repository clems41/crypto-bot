package trading

import "crypto-bot/internal/constant/currencyConst"

const (
	nbPreviousValues                            = 10
	delayBeforeNewAlgoApplicationInMilliseconds = 100
	gainToClosePositionInPercent                = 0.5 // meaning 0.5%
	nbOfIncreasingValueToOpenPosition           = 6
	maxOpenedPositions                          = 5
)

var (
	currenciesToTrade = []string{
		currencyConst.BtcEurPair,
	}
)
