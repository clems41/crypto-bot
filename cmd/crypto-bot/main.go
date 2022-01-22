package main

import (
	"crypto-bot/internal/domain/trading"
	"crypto-bot/internal/repository/localRepository"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/internal/service/tradingPlatform/krakenApi"
	"crypto-bot/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle signal interruption
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt, os.Kill)

	// Instantiate useful services and repositories
	//mockTrading, err := tradingPlatformMock.New()
	krakenApi, err := krakenApi.New()
	if err != nil {
		logger.Fatal(err)
	}
	priceRepo, err := localRepository.NewPriceRepository()
	if err != nil {
		logger.Fatal(err)
	}
	positionRepo, err := localRepository.NewPositionRepository()
	if err != nil {
		logger.Fatal(err)
	}

	// Run trading algorithm
	tradingService, err := trading.NewService(
		nil, // used default config or from env
		[]tradingPlatform.Api{
			krakenApi,
		},
		priceRepo,
		positionRepo)
	if err != nil {
		logger.Fatal(err)
	}
	go func() {
		err = tradingService.Start()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// This code below will be executed only when signal is received (shutdown, cancel, kill, etc...).
	// Add all closing methods here.
	<-stopSignal
	err = tradingService.Stop()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Infof("Crypto bot has been stopped !")
	}
}
