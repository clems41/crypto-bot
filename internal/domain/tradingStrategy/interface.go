package tradingStrategy

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"time"
)

type Algo interface {
	// ShouldAddOrder determines if order should be open. If returned order is nil, order should not be open
	ShouldAddOrder(form ShouldAddOrderForm) (shouldOpen bool, order model.Order, err error)
	// PricesNeeded return number of prices needed by ShouldAddOrder to work.
	// If ShouldAddOrder got fewer prices than value returned by PricesNeeded, ShouldAddOrder will return an error.
	PricesNeeded() (numberOfPrices int)
	// MaxOpenedOrdersByPair return limit of opened order by pair
	MaxOpenedOrdersByPair() (maxOpenedOrdersByPair int)
	// DelayBetweenEachRun return duration to wait between each algorithm execution
	DelayBetweenEachRun() (delay time.Duration)
	// PairsToTradeByPlatform return pairs that should be trades depending on platform
	PairsToTradeByPlatform() (pairsByPlatform map[string][]tradingConst.Pair)
}
