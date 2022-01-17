package trading

import (
	"crypto-bot/internal/service/tradingPlatform"
)

var _ Service = (*service)(nil)

type Service interface {
	Start() (err error)
	Stop() (err error)
}

type service struct {
	cryptoAPI tradingPlatform.Api
}

func NewService(cryptoAPI tradingPlatform.Api) Service {
	return &service{
		cryptoAPI: cryptoAPI,
	}
}

func (svc *service) Start() (err error) {
	return
}

func (svc *service) Stop() (err error) {
	return
}
