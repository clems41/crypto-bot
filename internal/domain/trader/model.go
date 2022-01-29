package trader

type TradeInfo struct {
	PairsToTrade      []string
	NbOpenOrderByPair map[string]int
}
