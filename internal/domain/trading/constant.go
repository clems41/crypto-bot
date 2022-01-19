package trading

import "crypto-bot/internal/constant/tradingConst"

const (
	delayBeforeNewAlgoApplicationInMilliseconds = 5000
	resultToClosePositionInPercent              = 0.15 // if equal to 0.1 and amount = 100, result should be at least 0.1 to send closing order
	nbOfIncreasingValueToOpenPosition           = 6
	maxOpenedPositionsByCurrency                = 1
	maxOpenedPositionsByPair                    = 1
	minimumAmountToOpenPosition                 = 10 // position will be open only if balance is at least equal to this minimum
)

var (
	pairsToTradeByPlatform = map[string][]string{
		tradingConst.KrakenPlatform: {
			tradingConst.BtcEurPair,
			tradingConst.DashEurPair,
		},
	}
)
