package maxOrMinAlgo

import (
	"crypto-bot/internal/model"
	"crypto-bot/internal/service/tradingStrategy"
)

var _ tradingStrategy.Algo = (*algorithm)(nil)

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

func (algo *algorithm) ShouldAddOrder(form tradingStrategy.ShouldAddOrderForm) (ok bool, order *model.Order, err error) {
	return
}
