package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/utils/envUtils"
	krakenClient "github.com/beldur/kraken-go-api-client"
)

var _ tradingPlatform.Api = (*krakenApi)(nil)

type krakenApi struct {
	client *krakenClient.KrakenAPI
}

func New() (api *krakenApi, err error) {
	apiKey, err := envUtils.GetFromEnvOrError(envKrakenApiKey)
	if err != nil {
		return nil, err
	}
	apiSecret, err := envUtils.GetFromEnvOrError(envKrakenApiSecret)
	if err != nil {
		return nil, err
	}
	client := krakenClient.New(apiKey, apiSecret)

	return &krakenApi{
		client: client,
	}, nil
}

func (api *krakenApi) Name() (name string) {
	return tradingConst.KrakenMockPlatform
}

func (api *krakenApi) AddOrder(order *model.Order) (err error) {
	err = order.Validate()
	if err != nil {
		return
	}

	// Send request to kraken but with validate=true (order will not be sent, but fields will be validated)
	return
}

func (api *krakenApi) CancelOrder(orderID string) (err error) {
	return
}

func (api *krakenApi) CancelAllOrders() (view tradingPlatform.CancelAllOrdersView, err error) {
	return
}

func (api *krakenApi) GetPrices(form tradingPlatform.GetPricesForm) (prices []*model.Price, err error) {
	return
}

func (api *krakenApi) GetIndexPrices(pairs ...string) (price map[string]*model.Price, err error) {
	return
}

func (api *krakenApi) GetOpenOrders() (orders []*model.Order, err error) {
	return
}

func (api *krakenApi) GetAllOrders() (orders []*model.Order, err error) {
	return
}

func (api *krakenApi) GetBalance() (balance *model.Balance, err error) {
	return
}
