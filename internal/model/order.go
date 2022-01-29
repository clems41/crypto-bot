package model

import (
	"crypto-bot/internal/constant/timeConst"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Order struct {
	ID                  string
	Date                time.Time
	Pair                string
	Side                string  `validate:"oneof=buy sell"`                           // buy or sell
	Volume              float64 `validate:"gte=0"`                                    // quantity of currency to buy/sell, can be 0, will be filled by trading platform
	Type                string  `validate:"oneof=market limit stop-loss take-profit"` // market, limit, stop-loss, take-profit
	Price               float64 `validate:"gt=0"`                                     // price of traded pair
	Amount              float64 `validate:"gt=0"`                                     // amount of initial currency to spend to buy another one
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
	return
}

func (order Order) String() (str string) {
	return fmt.Sprintf("{ %s : %s %s at %s with volume=%f, price=%f, amount=%f, status=%s, fees=%f, closeType=%s, closePrice=%f}",
		order.PlatformName, order.Side, order.Pair, order.Date.Format(timeConst.DefaultFormatTimeLayout), order.Volume,
		order.Price, order.Amount, order.Status, order.Fees, order.CloseConditionType, order.CloseConditionPrice)
}
