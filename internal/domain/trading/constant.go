package trading

import "crypto-bot/internal/constant/currencyConst"

const (
	delayBeforeNewAlgoApplicationInMilliseconds = 5000
	resultToClosePositionInPercent              = 0.15 // meaning 0.05% of gain so if equal to 0.1 and amount = 100, result should be at least 0.1 to send closing order
	nbOfIncreasingValueToOpenPosition           = 6
	maxOpenedPositions                          = 5
	balanceRatioToInvest                        = 0.75 // meaning that if balance = 100, 75 will be invested
	minimumToOpenPosition                       = 10   // position will be open only if balance is at least equal to this minimum
)

var (
	pairsToTrade = []string{
		currencyConst.BtcEurPair,
		currencyConst.DashEurPair,
	}
)
