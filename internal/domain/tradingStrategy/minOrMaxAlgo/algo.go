package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	tradingStrategy2 "crypto-bot/internal/domain/tradingStrategy"
	"crypto-bot/internal/model"
	"fmt"
)

var _ tradingStrategy2.Algo = (*algorithm)(nil)

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

func (algo *algorithm) ShouldAddOrder(form tradingStrategy2.ShouldAddOrderForm) (shouldOpen bool, order model.Order, err error) {
	if len(form.PriceHistory) < algo.PricesNeeded() {
		err = fmt.Errorf("not enough prices to take a decision, got %d but need %d",
			len(form.PriceHistory), algo.PricesNeeded())
	}

	// do calculation only on last NumberOfPreviousPricesToCompare prices
	pricesHistory := form.PriceHistory[len(form.PriceHistory)-algo.config.NumberOfPreviousPricesToCompare:]

	// find min from price history
	minimumAsk := pricesHistory[0].Ask
	for _, price := range pricesHistory {
		if price.Ask < minimumAsk {
			minimumAsk = price.Ask
		}
	}

	// fill order if conditions are ok
	// TODO fill amount and volume
	if form.IndexPrice.Ask <= minimumAsk {
		price := form.IndexPrice.Ask * (1 - algo.config.PercentPriceBelowToBuy/100)
		closeConditionPrice := price * (1 + algo.config.MinimumResultInPercentToClosePosition/100)
		order = model.Order{
			Pair:                form.Pair,
			Side:                tradingConst.BuySideOrder,
			Type:                tradingConst.LimitOrderType,
			Price:               price,
			CloseConditionType:  tradingConst.LimitCloseConditionType,
			CloseConditionPrice: closeConditionPrice,
		}
	}

	return
}

func (algo *algorithm) PricesNeeded() (numberOfPrices int) {
	return algo.config.NumberOfPreviousPricesToCompare
}

func (algo *algorithm) MaxOpenedOrdersByPair() (maxOpenedOrdersByPair int) {
	return algo.config.MaxOpenedOrdersByPair
}
