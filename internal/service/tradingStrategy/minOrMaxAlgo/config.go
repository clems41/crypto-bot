package minOrMaxAlgo

import (
	"crypto-bot/pkg/utils/envUtils"
	"strconv"
)

const (
	envNumberOfPreviousPricesToCompare = "NB_PRICES"
)

const (
	defaultNumberOfPreviousPricesToCompare = "60" // 1 price / min
)

type Config struct {
	// NumberOfPreviousPricesToCompare defines number of previous prices to use with opening position algorithm
	NumberOfPreviousPricesToCompare int
}

func GetConfigFromEnvOrDefault() (config *Config, err error) {
	numberOfPreviousPricesToCompareStr := envUtils.GetFromEnvOrDefault(envNumberOfPreviousPricesToCompare, defaultNumberOfPreviousPricesToCompare)
	numberOfPreviousPricesToCompare, err := strconv.Atoi(numberOfPreviousPricesToCompareStr)
	if err != nil {
		return
	}
	config = &Config{
		NumberOfPreviousPricesToCompare: numberOfPreviousPricesToCompare,
	}
	return
}
