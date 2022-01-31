package minOrMaxAlgo

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/domain/tradingStrategy"
	"crypto-bot/internal/model"
	"crypto-bot/pkg/utils/tradingUtils"
	"fmt"
	"github.com/pkg/errors"
)

var _ tradingStrategy.Algo = (*algorithm)(nil)

type algorithm struct {
	config *Config
}

// New instantiate new algorithm service. You can specify config parameters using config argument or environment variables.
func New(customConfig *Config) (algo *algorithm, err error) {
	var config *Config
	if customConfig != nil {
		config = customConfig
	} else {
		config, err = GetConfigFromEnvOrDefault()
		if err != nil {
			return
		}
	}

	algo = &algorithm{
		config: config,
	}
	return
}

func (algo *algorithm) ShouldAddOrder(form tradingStrategy.ShouldAddOrderForm) (shouldOpen bool, order model.Order, err error) {
	if len(form.PriceHistory) < algo.PricesNeeded() {
		err = fmt.Errorf("not enough prices to take a decision, got %d but need %d",
			len(form.PriceHistory), algo.PricesNeeded())
		return
	}

	// do calculation only on last NumberOfPreviousPricesToCompare prices
	pricesHistory := form.PriceHistory[len(form.PriceHistory)-algo.config.NumberOfPreviousPricesToCompare:]

	// find min and max from price history
	minimumAsk := pricesHistory[0].Ask
	maximumAsk := pricesHistory[0].Ask
	for _, price := range pricesHistory {
		if price.Ask < minimumAsk {
			minimumAsk = price.Ask
		}
		if price.Ask > maximumAsk {
			maximumAsk = price.Ask
		}
	}

	// fill order if conditions are ok
	//if form.IndexPrice.Ask <= minimumAsk {
	if form.IndexPrice.Ask >= maximumAsk {
		var amount float64
		amount, err = algo.getAmountToBuy(form)
		if err != nil {
			return false, order, errors.WithStack(err)
		}
		if amount < algo.config.MinimumAmount {
			return
		}
		price := form.IndexPrice.Ask * (1 - algo.config.PercentPriceBelowToBuy/100)
		price, err = tradingUtils.RemovePriceDecimal(price, form.PairToTrade)
		if err != nil {
			return false, order, errors.WithStack(err)
		}
		volume := amount / price
		closeConditionPrice := price * (1 + algo.config.MinimumResultInPercentToClosePosition/100)
		closeConditionPrice, err = tradingUtils.RemovePriceDecimal(closeConditionPrice, form.PairToTrade)
		if err != nil {
			return false, order, errors.WithStack(err)
		}
		order = model.Order{
			Pair:                form.PairToTrade,
			Side:                tradingConst.BuySideOrder,
			Type:                tradingConst.LimitOrderType,
			Price:               price,
			Volume:              volume,
			Amount:              amount,
			CloseConditionType:  tradingConst.LimitCloseConditionType,
			CloseConditionPrice: closeConditionPrice,
		}
		shouldOpen = true
	}

	return
}

func (algo *algorithm) PricesNeeded() (numberOfPrices int) {
	return algo.config.NumberOfPreviousPricesToCompare
}

func (algo *algorithm) MaxOpenedOrdersByPair() (maxOpenedOrdersByPair int) {
	return algo.config.MaxOpenedOrdersByPair
}

func (algo *algorithm) getAmountToBuy(form tradingStrategy.ShouldAddOrderForm) (amount float64, err error) {
	// find currency needed and related balance
	currencyNeededToBuy, ok := tradingUtils.CurrencyNeededToTradePair(form.PairToTrade, tradingConst.BuySideOrder)
	if !ok {
		return amount, fmt.Errorf("cannot find currency to buy %s", form.PairToTrade)
	}
	balance, ok := form.CurrentBalance[currencyNeededToBuy]
	if !ok {
		return amount, fmt.Errorf("cannot find balance for currency %s", currencyNeededToBuy)
	}
	if balance <= 0 {
		return
	}

	// count all pair that are using this currency
	var nbPairUsingCurrency int
	for _, pair := range form.AllPairsTraded {
		var currency string
		currency, ok = tradingUtils.CurrencyNeededToTradePair(pair, tradingConst.BuySideOrder)
		if !ok {
			return amount, fmt.Errorf("cannot find currency to buy %s", form.PairToTrade)
		}
		if currency == currencyNeededToBuy {
			nbPairUsingCurrency++
		}
	}

	// define share (number of part to divide balance)
	share := nbPairUsingCurrency * algo.MaxOpenedOrdersByPair()

	// count number of current open orders that already used currency balance
	var nbOpenOrderForPair, nbOpenOrderThatAlreadyUsedCurrency int
	for _, order := range form.OpenOrders {
		if order.Pair == form.PairToTrade && order.Status == tradingConst.OpenOrderStatus {
			nbOpenOrderForPair++
		}
		var currency string
		currency, ok = tradingUtils.CurrencyNeededToTradePair(order.Pair, tradingConst.BuySideOrder)
		if !ok {
			return amount, fmt.Errorf("cannot find currency to buy %s", form.PairToTrade)
		}
		if currency == currencyNeededToBuy && order.Side == tradingConst.SellSideOrder {
			nbOpenOrderThatAlreadyUsedCurrency++
		}
	}

	// order should not be open if MaxOpenedOrdersByPair has been reached
	if nbOpenOrderForPair >= algo.MaxOpenedOrdersByPair() {
		return
	}

	if share != nbOpenOrderThatAlreadyUsedCurrency {
		amount = balance / float64(share-nbOpenOrderThatAlreadyUsedCurrency)
	}

	return
}
