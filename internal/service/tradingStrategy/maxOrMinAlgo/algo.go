package maxOrMinAlgo

import (
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

func (algo *algorithm) ShouldAddOrder() (ok bool, err error) {
	return
}

func (algo *algorithm) Config() (config *tradingStrategy.Config, err error) {
	config = &tradingStrategy.Config{
		NumberPreviousPricesNeeded: algo.config.NumberOfPreviousPricesToCompare,
	}
	return
}
