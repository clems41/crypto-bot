package mathUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveNDecimal(t *testing.T) {
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
		actual, err := RemoveNDecimal(args.Value, args.NbDecimals)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}
