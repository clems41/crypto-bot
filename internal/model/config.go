package model

import "github.com/go-playground/validator/v10"

type Config struct {
	Name       string `validate:"required"`
	Parameters string `validate:"required"`
}

func (cfg *Config) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return
	}
	return
}
