package model

import (
	"crypto-bot/internal/constant/timeConst"
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/pkg/utils/tradingUtils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Order struct {
	ID                  string    `validate:"required"`
	OpenTime            time.Time `validate:"required"`
	CloseTime           time.Time
	Pair                tradingConst.Pair        `validate:"required"`
	Side                tradingConst.OrderSide   `validate:"oneof=buy sell"`                           // buy or sell
	Volume              float64                  `validate:"gte=0"`                                    // quantity of currency to buy/sell, can be 0, will be filled by trading platform
	Type                tradingConst.OrderType   `validate:"oneof=market limit stop-loss take-profit"` // market, limit, stop-loss, take-profit
	Price               float64                  `validate:"gte=0"`                                    // price of traded pair
	Amount              float64                  `validate:"gte=0"`                                    // amount of initial currency to spend to buy another one
	Leverage            int                      `validate:"gte=0"`                                    // effet de levier x1, x2 ,x3, etc...
	CloseConditionType  tradingConst.OrderType   `validate:"oneof=none limit stop-loss take-profit"`   // condition to create an opposite order when the first one is completed : limit, stop-loss, take-profit
	CloseConditionPrice float64                  `validate:"gte=0"`                                    // price that opposite order should get before executing order
	Fees                float64                  `validate:"gte=0"`                                    // fees taken by trading platform
	Status              tradingConst.OrderStatus `validate:"oneof=open close cancel"`                  // order status
	PlatformName        string                   `validate:"required"`                                 // name of platform
}

func (order Order) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		return
	}

	// check status
	if order.Status != tradingConst.Cancel && order.Price == 0 {
		return fmt.Errorf("price should not be 0 if status is not cancel")
	}

	// check pair price decimals
	expectedPrice, err := tradingUtils.RemovePriceDecimal(order.Price, order.Pair)
	if err != nil {
		return
	}
	if expectedPrice != order.Price {
		return fmt.Errorf("price should be rounded to %f but it is %f", expectedPrice, order.Price)
	}

	// check amount decimals
	currency, ok := tradingUtils.CurrencyNeededToTradePair(order.Pair, tradingConst.Buy)
	if !ok {
		return fmt.Errorf("currency needed cannot be found for pair %s", order.Pair)
	}
	expectedAmount, err := tradingUtils.RemoveAmountDecimalWithFloorRounding(order.Amount, currency)
	if err != nil {
		return
	}
	if expectedAmount != order.Amount {
		return fmt.Errorf("amount should be rounded to %f but it is %f", expectedAmount, order.Amount)
	}

	return
}

func (order Order) String() (str string) {
	return fmt.Sprintf("{ %s : %s %s at (open=%s close=%s) with volume=%f, price=%f, amount=%f, status=%s, fees=%f, closeType=%s, closePrice=%f}",
		order.PlatformName, order.Side, order.Pair, order.OpenTime.Format(timeConst.DefaultFormatTimeLayout),
		order.CloseTime.Format(timeConst.DefaultFormatTimeLayout), order.Volume, order.Price, order.Amount,
		order.Status, order.Fees, order.CloseConditionType, order.CloseConditionPrice)
}
