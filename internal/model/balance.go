package model

import (
	"github.com/go-playground/validator/v10"
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
