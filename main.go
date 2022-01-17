package main

import (
	"crypto-bot/internal/domain/trading"
	"crypto-bot/internal/service/tradingPlatform/tradingPlatformMock"
	"crypto-bot/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle signal interruption
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt, os.Kill)

	tradingApi := tradingPlatformMock.New()
	tradingService := trading.NewService(tradingApi)
	go func() {
		err := tradingService.Start()
		if err != nil {
			logger.Fatalf("Cannot start crypto bot : %v", err)
		} else {
			logger.Infof("Crypto bot is running...")
		}
	}()

	// This code below will be executed only when signal is received (shutdown, cancel, kill, etc...).
	// Add all closing methods here.
	<-stopSignal
	err := tradingService.Stop()
	if err != nil {
		logger.Fatalf("Cannot stop crypto bot : %v", err)
	} else {
		logger.Infof("Crypto bot has been stopped !")
	}
}
