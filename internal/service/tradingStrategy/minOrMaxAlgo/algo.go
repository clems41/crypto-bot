package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/service/tradingStrategy"
	"fmt"
)

var _ tradingStrategy.Algo = (*algorithm)(nil)

type algorithm struct {
	config *Config
}

func New() (algo *algorithm, err error) {
	config, err := GetConfigFromEnvOrDefault()
	if err != nil {
		return
	}

	algo = &algorithm{
		config: config,
	}
	return
}

func (algo *algorithm) ShouldAddOrder(form tradingStrategy.ShouldAddOrderForm) (view tradingStrategy.ShouldAddOrderView, err error) {
	if len(form.PriceHistory) < algo.config.NumberOfPreviousPricesToCompare {
		err = fmt.Errorf("not enough prices to take a decision, got %d but need %d",
			len(form.PriceHistory), algo.config.NumberOfPreviousPricesToCompare)
	}

	// do calculation only on last NumberOfPreviousPricesToCompare prices
	pricesHistory := form.PriceHistory[len(form.PriceHistory)-algo.config.NumberOfPreviousPricesToCompare:]

	// find min and max from price history
	minimumAsk := pricesHistory[0].Ask
	maximumBid := pricesHistory[0].Bid
	for _, price := range pricesHistory {
		if price.Ask < minimumAsk {
			minimumAsk = price.Ask
		}
		if price.Bid > maximumBid {
			maximumBid = price.Bid
		}
	}

	// fill order if conditions are ok
	if form.IndexPrice.Ask <= minimumAsk {
		view = tradingStrategy.ShouldAddOrderView{
			ShouldAddOrder: true,
			Side:           tradingConst.BuySideOrder,
			Type:           tradingConst.LimitOrderType,
			Price:          form.IndexPrice.Ask,
		}
	} else if form.IndexPrice.Bid >= maximumBid {
		view = tradingStrategy.ShouldAddOrderView{
			ShouldAddOrder: true,
			Side:           tradingConst.SellSideOrder,
			Type:           tradingConst.LimitOrderType,
			Price:          form.IndexPrice.Bid,
		}
	}

	return
}
