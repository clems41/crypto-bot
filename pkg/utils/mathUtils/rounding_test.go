package mathUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoundingAfterNDecimal(t *testing.T) {
	type arg struct {
		Value      float64
		NbDecimals int
	}
	tests := map[arg]float64{
		{
			Value:      10.568,
			NbDecimals: 1,
		}: 10.6,
		{
			Value:      10.568,
			NbDecimals: 3,
		}: 10.568,
		{
			Value:      10.568,
			NbDecimals: 4,
		}: 10.568,
		{
			Value:      10,
			NbDecimals: 4,
		}: 10,
	}
	for args, expected := range tests {
		actual, err := RoundingAfterNDecimal(args.Value, args.NbDecimals)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}

func TestFloorRoundingAfterNDecimal(t *testing.T) {
	type arg struct {
		Value      float64
		NbDecimals int
	}
	tests := map[arg]float64{
		{
			Value:      10.568,
			NbDecimals: 1,
		}: 10.5,
		{
			Value:      10.5682,
			NbDecimals: 3,
		}: 10.568,
		{
			Value:      10.5689,
			NbDecimals: 3,
		}: 10.568,
		{
			Value:      10.568,
			NbDecimals: 4,
		}: 10.568,
		{
			Value:      10,
			NbDecimals: 4,
		}: 10,
	}
	for args, expected := range tests {
		actual, err := FloorRoundingAfterNDecimal(args.Value, args.NbDecimals)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}

func TestCeilRoundingAfterNDecimal(t *testing.T) {
	type arg struct {
		Value      float64
		NbDecimals int
	}
	tests := map[arg]float64{
		{
			Value:      10.568,
			NbDecimals: 1,
		}: 10.6,
		{
			Value:      10.5682,
			NbDecimals: 3,
		}: 10.569,
		{
			Value:      10.5689,
			NbDecimals: 3,
		}: 10.569,
		{
			Value:      10.568,
			NbDecimals: 4,
		}: 10.568,
		{
			Value:      10,
			NbDecimals: 4,
		}: 10,
	}
	for args, expected := range tests {
		actual, err := CeilRoundingAfterNDecimal(args.Value, args.NbDecimals)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}
