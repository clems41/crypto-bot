package krakenApiMock

import (
	"crypto-bot/external/service/tradingPlatform"
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
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
	closeOrderType, ok := typeConverter[order.CloseConditionType]
	if !ok {
		return fmt.Errorf("cannot find close order type for %s", order.CloseConditionType)
	}
	orderParameters := map[string]string{
		priceParameter:    fmt.Sprintf("%f", order.Price),
		validateParameter: "true",
	}
	if closeOrderType != "" {
		orderParameters[closeOrderTypeParameter] = closeOrderType
		orderParameters[closePriceParameter] = fmt.Sprintf("%f", order.CloseConditionPrice)
	}
	response, err := api.client.AddOrder(pair, side, orderType, fmt.Sprintf("%f", order.Volume), orderParameters)
	if err != nil {
		return
	}
	logger.Infof("Order sent to kraken api : %+v", response.Description)

	// fill order as mock
	fees := order.Amount * takerFees / 100
	order.ID = uuid.New().String()
	order.Fees = fees
	order.Status = tradingConst.OpenOrderStatus
	order.Date = time.Now()

	// add orders in memory
	api.orders = append(api.orders, order)
	return
}

func (api *krakenApi) GetIndexPrices(pairs ...string) (prices []model.Price, err error) {
	if len(pairs) == 0 {
		return
	}
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

	err = api.updateOrdersBasedOnPrice(prices)
	if err != nil {
		return
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

// updateOrdersBasedOnPrice will update orders like a real platform will do.
// Close open orders if price has reached the one fixed in order.
// Open new order if close condition is defined.
func (api *krakenApi) updateOrdersBasedOnPrice(prices []model.Price) (err error) {
	// map prices by pair
	indexPricesByPair := make(map[string]model.Price)
	for _, price := range prices {
		indexPricesByPair[price.Pair] = price
	}

	// update open orders and create new one if conditions are reached
	for orderIdx, order := range api.orders {
		if order.Status == tradingConst.OpenOrderStatus {
			// find if order should be closed
			var shouldClose bool
			if order.Side == tradingConst.BuySideOrder {
				if order.Price >= indexPricesByPair[order.Pair].Ask {
					shouldClose = true
				}
			} else {
				if order.Price <= indexPricesByPair[order.Pair].Bid {
					shouldClose = true
				}
			}

			// if yes, close it and update balance
			if shouldClose {
				err = api.closeOrder(order)
				if err != nil {
					return
				}
				api.orders[orderIdx] = order

				// if order close condition are specified, create new order based on it
				if order.CloseConditionType != "" && order.CloseConditionType != tradingConst.NoneCloseConditionType {
					if order.CloseConditionType == tradingConst.TakeProfitCloseConditionType ||
						order.CloseConditionType == tradingConst.LimitCloseConditionType {
						var closeOrderSide string
						if order.Side == tradingConst.BuySideOrder {
							closeOrderSide = tradingConst.SellSideOrder
						} else {
							closeOrderSide = tradingConst.BuySideOrder
						}
						newOrder := model.Order{
							ID:                 uuid.New().String(),
							Date:               time.Now(),
							Pair:               order.Pair,
							Side:               closeOrderSide,
							Type:               order.CloseConditionType,
							Price:              order.CloseConditionPrice,
							PlatformName:       api.Name(),
							CloseConditionType: tradingConst.NoneCloseConditionType,
						}
						if order.Side == tradingConst.BuySideOrder {
							newOrder.Volume = order.Volume
							newOrder.Amount = newOrder.Price * newOrder.Volume
						} else {
							newOrder.Amount = order.Amount
							newOrder.Volume = newOrder.Amount / newOrder.Price
						}
						err = api.AddOrder(&newOrder)
						if err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}

func (api *krakenApi) closeOrder(order *model.Order) (err error) {
	order.Status = tradingConst.CloseOrderStatus
	var currencyNeeded, currencyGot string
	var ok bool
	currencyNeeded, ok = tradingUtils.CurrencyNeededToTradePair(order.Pair, order.Side)
	if !ok {
		return fmt.Errorf("cannot find currency needed for pair %s ans side %s", order.Pair, order.Side)
	}
	currencyGot, ok = tradingUtils.CurrencyGotAfterTradingPair(order.Pair, order.Side)
	if !ok {
		return fmt.Errorf("cannot find currency got for pair %s ans side %s", order.Pair, order.Side)
	}
	if order.Side == tradingConst.BuySideOrder {
		order.Volume *= 1 - takerFees/100
		api.balanceByCurrencyMock[currencyNeeded] -= order.Amount
		api.balanceByCurrencyMock[currencyGot] += order.Volume
	} else {
		order.Amount *= 1 - takerFees/100
		api.balanceByCurrencyMock[currencyNeeded] -= order.Volume
		api.balanceByCurrencyMock[currencyGot] += order.Amount
	}
	logger.Infof("Closing order %s", order.ID)
	return
}
