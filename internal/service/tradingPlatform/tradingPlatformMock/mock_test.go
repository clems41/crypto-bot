package tradingPlatformMock

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMock_GetWalletBalance(t *testing.T) {
	testMock := New()
	walletBalance, err := testMock.GetWalletBalance()
	require.NoError(t, err)
	require.Equal(t, initWalletBalance, walletBalance)
}
