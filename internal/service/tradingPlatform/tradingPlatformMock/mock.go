package tradingPlatformMock

import "crypto-bot/internal/service/tradingPlatform"

var _ tradingPlatform.Api = (*mock)(nil)

type mock struct {
}

func New() *mock {
	return &mock{}
}

func (mock *mock) GetWalletBalance() (amount int, err error) {
	return
}
