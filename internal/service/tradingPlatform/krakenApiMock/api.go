package krakenApiMock

import (
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"crypto-bot/internal/service/tradingPlatform"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
	"crypto-bot/pkg/utils/tradingUtils"
	"fmt"
	krakenClient "github.com/beldur/kraken-go-api-client"
	"github.com/google/uuid"
	"reflect"
	"strconv"
	"time"
)

var _ tradingPlatform.Api = (*krakenApi)(nil)

type krakenApi struct {
	client                *krakenClient.KrakenAPI
	balanceByCurrencyMock map[string]float64
	orders                []*model.Order
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
		client:                client,
		balanceByCurrencyMock: initialBalance,
	}, nil
}

func (api *krakenApi) Name() (name string) {
	return tradingConst.KrakenMockPlatform
}

func (api *krakenApi) AddOrder(order *model.Order) (err error) {
	// Send request to kraken but with validate=true (order will not be sent, but fields will be validated)
	pair, err := GetKrakenPair(order.Pair)
	if err != nil {
		return
	}
	side, ok := sideConverter[order.Side]
	if !ok {
		return fmt.Errorf("cannot find side for %s", order.Side)
	}
	orderType, ok := typeConverter[order.Type]
	if !ok {
		return fmt.Errorf("cannot find order type for %s", order.Type)
	}
	response, err := api.client.AddOrder(pair, side, orderType, fmt.Sprintf("%f", order.Volume), map[string]string{
		priceParameter:    fmt.Sprintf("%f", order.Price),
		validateParameter: "true",
	})
	if err != nil {
		return
	}

	// TODO fill order from platform orders info
	logger.Infof("Orders %v has been added to %s : %s", api.Name(), response.TransactionIds, response.Description)

	// fill order as mock
	fees := order.Amount * takerFees / 100
	volumeAfterFees := order.Volume * (1 - takerFees/100)
	amountAfterFees := order.Amount * (1 - takerFees/100)
	order.ID = uuid.New().String()
	order.Fees = fees
	order.Status = tradingConst.CloseOrderStatus

	// update balance mock
	currencyNeeded, ok := tradingUtils.CurrencyNeededToTradePair(order.Pair, order.Side)
	if !ok {
		return fmt.Errorf("cannot find currency for pair %s", order.Pair)
	}
	currencyGot, ok := tradingUtils.CurrencyGotAfterTradingPair(order.Pair, order.Side)
	if !ok {
		return fmt.Errorf("cannot find currency for pair %s", order.Pair)
	}
	if order.Side == tradingConst.BuySideOrder {
		api.balanceByCurrencyMock[currencyNeeded] -= order.Amount
		api.balanceByCurrencyMock[currencyGot] += volumeAfterFees
	} else {
		api.balanceByCurrencyMock[currencyNeeded] -= order.Volume
		api.balanceByCurrencyMock[currencyGot] += amountAfterFees
	}

	// add orders in memory
	api.orders = append(api.orders, order)
	return
}

// CancelAllOrders with mock all orders will be instantaneously open, executed then close.
// There will be no orders left to cancel
func (api *krakenApi) CancelAllOrders() (view tradingPlatform.CancelAllOrdersView, err error) {
	return
}

func (api *krakenApi) GetIndexPrices(pairs ...string) (prices []model.Price, err error) {
	// get response form kraken api
	var krakenPairs []string
	for _, pair := range pairs {
		var krakenPair string
		krakenPair, err = GetKrakenPair(pair)
		if err != nil {
			return
		}
		krakenPairs = append(krakenPairs, krakenPair)
	}
	response, err := api.client.Ticker(krakenPairs...)
	if err != nil {
		return
	}

	// search required value from response
	value := reflect.ValueOf(response).Elem()
	typeOfResponse := value.Type()
	for i := 0; i < value.NumField(); i++ {
		pairTickerInfoInterface := value.Field(i).Interface()
		pairTickerInfo, ok := pairTickerInfoInterface.(krakenClient.PairTickerInfo)
		if ok && len(pairTickerInfo.Ask) > 0 && len(pairTickerInfo.Bid) > 0 {
			var pair string
			pair, err = GetProjectPair(typeOfResponse.Field(i).Name)
			if err != nil {
				return
			}
			var askPrice, bidPrice float64
			askPrice, err = strconv.ParseFloat(pairTickerInfo.Ask[0], 64)
			if err != nil {
				return
			}
			bidPrice, err = strconv.ParseFloat(pairTickerInfo.Bid[0], 64)
			if err != nil {
				return
			}
			prices = append(prices, model.Price{
				Date:         time.Now(),
				PlatformName: api.Name(),
				Pair:         pair,
				Ask:          askPrice,
				Bid:          bidPrice,
			})
		}
	}

	return
}

func (api *krakenApi) GetOpenOrders() (orders []model.Order, err error) {
	for _, order := range api.orders {
		if order.Status == tradingConst.OpenOrderStatus {
			if order != nil {
				orders = append(orders, *order)
			}
		}
	}
	return
}

func (api *krakenApi) GetAllOrders() (orders []model.Order, err error) {
	for _, order := range api.orders {
		if order != nil {
			orders = append(orders, *order)
		}
	}
	return
}

func (api *krakenApi) RefreshBalance(balance *model.Balance) (err error) {
	balance.ValueByCurrency = api.balanceByCurrencyMock
	balance.UpdatedAt = time.Now()
	balance.PlatformName = api.Name()
	return
}
