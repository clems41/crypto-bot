package main

import (
	"context"
	"crypto-bot/internal/domain/trader"
	"crypto-bot/internal/repository/googleSheetRepository"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/internal/service/tradingPlatform/krakenApiMock"
	"crypto-bot/internal/service/tradingStrategy/minOrMaxAlgo"
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
	krakenApi, err := krakenApiMock.New()
	if err != nil {
		logger.Fatal(err)
	}
	minOrMaxAlgo, err := minOrMaxAlgo.New()
	if err != nil {
		logger.Fatal(err)
	}
	googleSheetRepo, err := googleSheetRepository.New(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	// Run trading algorithm
	tradingService, err := trader.NewService(
		nil, // used default config or from env
		[]tradingPlatform.Api{
			krakenApi,
		},
		googleSheetRepo,
		minOrMaxAlgo)
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
