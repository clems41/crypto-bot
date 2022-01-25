package model

import (
	"crypto-bot/internal/constant/timeConst"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Price struct {
	Date         time.Time `validate:"required"`
	PlatformName string    `validate:"required"`
	Pair         string    `validate:"required"`
	Ask          float64   `validate:"gte=0"`
	Bid          float64   `validate:"gte=0"`
}

func (price *Price) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(price)
	if err != nil {
		return
	}
	return
}

func (price Price) String() (str string) {
	return fmt.Sprintf("%s at %s on %s --> ask=%f | bid=%f", price.Pair, price.Date.Format(timeConst.DefaultFormatTimeLayout),
		price.PlatformName, price.Ask, price.Bid)
}
