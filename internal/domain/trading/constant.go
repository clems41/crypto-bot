package trading

import "crypto-bot/internal/constant/tradingConst"

const (
	delayBeforeNewAlgoApplicationInMilliseconds = 5000
	minimumResultToClosePositionInPercent       = 0.2 // if equal to 0.1 and amount = 100, result should be at least 0.1 to send closing order
	nbPricesToCompare                           = 10  // specify number of last prices that should be taken into account to define if position should be open
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
