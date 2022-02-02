package model

import (
	"crypto-bot/internal/constant/timeConst"
	"crypto-bot/internal/constant/tradingConst"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Order struct {
	ID                  string    `validate:"required"`
	OpenTime            time.Time `validate:"required"`
	CloseTime           time.Time
	Pair                string  `validate:"required"`
	Side                string  `validate:"oneof=buy sell"`                           // buy or sell
	Volume              float64 `validate:"gte=0"`                                    // quantity of currency to buy/sell, can be 0, will be filled by trading platform
	Type                string  `validate:"oneof=market limit stop-loss take-profit"` // market, limit, stop-loss, take-profit
	Price               float64 `validate:"gte=0"`                                    // price of traded pair
	Amount              float64 `validate:"gte=0"`                                    // amount of initial currency to spend to buy another one
	Leverage            int     `validate:"gte=0"`                                    // effet de levier x1, x2 ,x3, etc...
	CloseConditionType  string  `validate:"oneof=none limit stop-loss take-profit"`   // condition to create an opposite order when the first one is completed : limit, stop-loss, take-profit
	CloseConditionPrice float64 `validate:"gte=0"`                                    // price that opposite order should get before executing order
	Fees                float64 `validate:"gte=0"`                                    // fees taken by trading platform
	Status              string  `validate:"oneof=open close cancel"`                  // order status
	PlatformName        string  `validate:"required"`                                 // name of platform
}

func (order Order) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		return
	}
	if order.Status == tradingConst.CloseOrderStatus && order.Price == 0 {
		return fmt.Errorf("price should not be 0 if status is close")

	}
	return
}

func (order Order) String() (str string) {
	return fmt.Sprintf("{ %s : %s %s at (open=%s close=%s) with volume=%f, price=%f, amount=%f, status=%s, fees=%f, closeType=%s, closePrice=%f}",
		order.PlatformName, order.Side, order.Pair, order.OpenTime.Format(timeConst.DefaultFormatTimeLayout),
		order.CloseTime.Format(timeConst.DefaultFormatTimeLayout), order.Volume, order.Price, order.Amount,
		order.Status, order.Fees, order.CloseConditionType, order.CloseConditionPrice)
}
