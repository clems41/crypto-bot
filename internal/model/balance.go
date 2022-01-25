package model

import (
	"crypto-bot/internal/constant/timeConst"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type Balance struct {
	PlatformName    string             `validate:"required"`
	ValueByCurrency map[string]float64 `validate:"required"`
	UpdatedAt       time.Time          `validate:"required"`
}

func (balance *Balance) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(balance)
	if err != nil {
		return
	}
	return
}

func (balance Balance) String() (str string) {
	str += fmt.Sprintf("Balance for %s at %s --> ", balance.PlatformName, balance.UpdatedAt.Format(timeConst.DefaultFormatTimeLayout))
	var balancesStr []string
	for currency, value := range balance.ValueByCurrency {
		balancesStr = append(balancesStr, fmt.Sprintf("%s : %f", currency, value))
	}
	str += strings.Join(balancesStr, " | ")
	return
}
