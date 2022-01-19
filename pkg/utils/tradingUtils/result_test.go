package tradingUtils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEstimateProfit(t *testing.T) {
	type arg struct {
		InitialBalance    float64
		FinalBalance      float64
		TradingDuration   time.Duration
		EstimatedDuration time.Duration
	}
	tests := map[arg]float64{
		{
			InitialBalance:    100,
			FinalBalance:      100,
			TradingDuration:   1 * time.Second,
			EstimatedDuration: 1 * time.Hour,
		}: 100,
		{
			InitialBalance:    100,
			FinalBalance:      100.1,
			TradingDuration:   1 * time.Minute,
			EstimatedDuration: 1 * time.Hour,
		}: 106.18,
		{
			InitialBalance:    100,
			FinalBalance:      99.99,
			TradingDuration:   1 * time.Hour,
			EstimatedDuration: 24 * time.Hour,
		}: 99.76,
		{
			InitialBalance:    100,
			FinalBalance:      100.1,
			TradingDuration:   1 * time.Hour,
			EstimatedDuration: 365 * 24 * time.Hour,
		}: 634627.25,
		{
			InitialBalance:    100,
			FinalBalance:      150,
			TradingDuration:   2 * time.Hour,
			EstimatedDuration: 5 * time.Hour,
		}: 275.57,
	}
	for args, expected := range tests {
		actual := EstimateProfit(args.InitialBalance, args.FinalBalance, args.TradingDuration, args.EstimatedDuration)
		require.Equal(t, fmt.Sprintf("%0.2f", expected), fmt.Sprintf("%0.2f", actual)) // compare only with 2 decimals
	}
}

func TestGetProfit(t *testing.T) {
	type arg struct {
		AskPrice float64
		BidPrice float64
		Amount   float64
	}
	tests := map[arg]float64{
		{
			AskPrice: 100,
			BidPrice: 100,
			Amount:   10,
		}: 10,
		{
			AskPrice: 100,
			BidPrice: 200,
			Amount:   10,
		}: 20,
		{
			AskPrice: 200,
			BidPrice: 100,
			Amount:   10,
		}: 5,
	}
	for args, expected := range tests {
		actual := GetProfit(args.AskPrice, args.BidPrice, args.Amount)
		require.Equal(t, expected, actual)
	}
}

func TestGetResult(t *testing.T) {
	type arg struct {
		AskPrice float64
		BidPrice float64
		Amount   float64
	}
	tests := map[arg]float64{
		{
			AskPrice: 100,
			BidPrice: 100,
			Amount:   10,
		}: 0.0,
		{
			AskPrice: 100,
			BidPrice: 200,
			Amount:   10,
		}: 10,
		{
			AskPrice: 200,
			BidPrice: 100,
			Amount:   10,
		}: -5,
	}
	for args, expected := range tests {
		actual := GetResult(args.AskPrice, args.BidPrice, args.Amount)
		require.Equal(t, expected, actual)
	}
}

func TestGetResultInPercent(t *testing.T) {
	type arg struct {
		AskPrice float64
		BidPrice float64
		Amount   float64
	}
	tests := map[arg]float64{
		{
			AskPrice: 100,
			BidPrice: 100,
			Amount:   10,
		}: 0.0,
		{
			AskPrice: 100,
			BidPrice: 200,
			Amount:   10,
		}: 100,
		{
			AskPrice: 200,
			BidPrice: 100,
			Amount:   10,
		}: -50,
	}
	for args, expected := range tests {
		actual := GetResultInPercent(args.AskPrice, args.BidPrice, args.Amount)
		require.Equal(t, expected, actual)
	}

}
