package trading

import "crypto-bot/internal/constant/currencyConst"

const (
	nbPreviousValues                            = 10
	delayBeforeNewAlgoApplicationInMilliseconds = 2000
	gainToClosePositionInPercent                = 0.05 // meaning 0.05% of gain so if equal to 0.1 and amount = 100, result should be at least 0.1 to send closing order
	nbOfIncreasingValueToOpenPosition           = 6
	maxOpenedPositions                          = 5
)

var (
	currenciesToTrade = []string{
		currencyConst.BtcEurPair,
	}
)
