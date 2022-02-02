package trader

import "crypto-bot/internal/constant/tradingConst"

type TradeInfo struct {
	PairsToTrade      []tradingConst.Pair
	NbOpenOrderByPair map[tradingConst.Pair]int
}
