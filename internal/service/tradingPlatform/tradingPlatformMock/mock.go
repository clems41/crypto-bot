package tradingPlatformMock

import "crypto-bot/internal/service/tradingPlatform"

var _ tradingPlatform.Api = (*mock)(nil)

type mock struct {
	balance float32
}

func New() *mock {
	return &mock{
		balance: initWalletBalance,
	}
}

func (mock *mock) GetWalletBalance() (amount float32, err error) {
	return mock.balance, nil
}
