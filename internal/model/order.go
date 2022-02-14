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
	Volume              float64                  `validate:"gt=0"`                                     // quantity of currency to buy/sell, can be 0, will be filled by trading platform
	Type                tradingConst.OrderType   `validate:"oneof=market limit stop-loss take-profit"` // market, limit, stop-loss, take-profit
	Price               float64                  `validate:"gt=0"`                                     // price of traded pair
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

	return
}

func (order Order) ValidateBeforeAdding() (err error) {
	// check pair
	if order.Pair == "" {
		return fmt.Errorf("pair should not be empty")
	}

	// check volume
	if order.Volume == 0 {
		return fmt.Errorf("volume should not be empty")
	}

	// check side
	if order.Side != tradingConst.Buy && order.Side != tradingConst.Sell {
		return fmt.Errorf("order side should be one of %v but not '%s'", []tradingConst.OrderSide{
			tradingConst.Buy,
			tradingConst.Sell,
		}, order.Side)
	}

	// check type
	if order.Type != tradingConst.None && order.Type != tradingConst.Market && order.Type != tradingConst.Limit &&
		order.Type != tradingConst.TakeProfit && order.Type != tradingConst.StopLoss {
		return fmt.Errorf("type should be one of %v but it is '%s'", []tradingConst.OrderType{
			tradingConst.None,
			tradingConst.Market,
			tradingConst.Limit,
			tradingConst.TakeProfit,
			tradingConst.StopLoss,
		}, order.Type)
	}

	// check close condition type
	if order.CloseConditionType != tradingConst.None && order.CloseConditionType != tradingConst.Market && order.CloseConditionType != tradingConst.Limit &&
		order.CloseConditionType != tradingConst.TakeProfit && order.CloseConditionType != tradingConst.StopLoss {
		return fmt.Errorf("close condition type should be one of %v but it is '%s'", []tradingConst.OrderType{
			tradingConst.None,
			tradingConst.Market,
			tradingConst.Limit,
			tradingConst.TakeProfit,
			tradingConst.StopLoss,
		}, order.CloseConditionType)
	}

	// check close condition price
	if order.CloseConditionType != "" && order.CloseConditionPrice == 0 {
		return fmt.Errorf("close condition price should not be empty with close condition %s", order.CloseConditionType)
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
