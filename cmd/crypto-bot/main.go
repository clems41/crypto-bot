package main

import (
	"crypto-bot/external/service/tradingPlatform"
	"crypto-bot/external/service/tradingPlatform/krakenApiMock"
	"crypto-bot/internal/domain/trader"
	"crypto-bot/internal/domain/tradingStrategy/minOrMaxAlgo"
	"crypto-bot/internal/repository/postgresqlRepository"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/sql/sqlConnection"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle signal interruption
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt, os.Kill)

	// Instantiate PostgresSQl database connection
	DB, err := sqlConnection.Open(nil)
	if err != nil {
		logger.Fatal(err)
	}

	/* Trading platforms */
	krakenApi, err := krakenApiMock.New()
	if err != nil {
		logger.Fatal(err)
	}

	/* Repository */
	/*	googleSheetRepo, err := googleSheetRepository.New(context.Background())
		if err != nil {
			logger.Fatal(err)
		}*/
	postgresqlRepo, err := postgresqlRepository.New(DB)
	if err != nil {
		logger.Fatal(err)
	}

	/* Algorithm */
	algo, err := minOrMaxAlgo.New()
	if err != nil {
		logger.Fatal(err)
	}

	// Run trading algorithm
	tradingService, err := trader.NewService(
		[]tradingPlatform.Api{
			krakenApi,
		},
		postgresqlRepo,
		algo)
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
	tradingService.Stop()
	logger.Infof("Crypto bot has been stopped !")
}
