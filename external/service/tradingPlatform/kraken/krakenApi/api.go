package krakenApi

import (
	"crypto-bot/external/service/tradingPlatform"
	"crypto-bot/external/service/tradingPlatform/kraken"
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
	"fmt"
	krakenClient "github.com/beldur/kraken-go-api-client"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var _ tradingPlatform.Api = (*krakenApi)(nil)

type krakenApi struct {
	client *krakenClient.KrakenAPI
}

func New() (api *krakenApi, err error) {
	apiKey, err := envUtils.GetFromEnvOrError(kraken.EnvKrakenApiKey)
	if err != nil {
		return nil, err
	}
	apiSecret, err := envUtils.GetFromEnvOrError(kraken.EnvKrakenApiSecret)
	if err != nil {
		return nil, err
	}
	client := krakenClient.New(apiKey, apiSecret)

	return &krakenApi{
		client: client,
	}, nil
}

func (api *krakenApi) Name() (name string) {
	return tradingConst.KrakenPlatform
}

func (api *krakenApi) AddOrder(order model.Order) (err error) {
	// Send request to kraken but with validate=true (order will not be sent, but fields will be validated)
	pair, err := kraken.GetKrakenPair(order.Pair)
	if err != nil {
		return errors.WithStack(err)
	}
	side, ok := kraken.SideConverter[order.Side]
	if !ok {
		return fmt.Errorf("cannot find side for %s", order.Side)
	}
	orderType, ok := kraken.TypeConverter[order.Type]
	if !ok {
		return fmt.Errorf("cannot find order type for %s", order.Type)
	}
	closeOrderType, ok := kraken.TypeConverter[order.CloseConditionType]
	if !ok {
		return fmt.Errorf("cannot find close order type for %s", order.CloseConditionType)
	}
	orderParameters := map[string]string{
		kraken.PriceParameter: fmt.Sprintf("%f", order.Price),
	}
	if closeOrderType != "" {
		orderParameters[kraken.CloseOrderTypeParameter] = closeOrderType
		orderParameters[kraken.ClosePriceParameter] = fmt.Sprintf("%f", order.CloseConditionPrice)
	}
	if !ok {
		return fmt.Errorf("cannot find close order type for %s", order.CloseConditionType)
	}
	_, err = api.client.AddOrder(pair, side, orderType, fmt.Sprintf("%f", order.Volume), orderParameters)
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (api *krakenApi) GetIndexPrices(pairs ...tradingConst.Pair) (prices []model.Price, err error) {
	// get response form kraken api
	var krakenPairs []string
	for _, pair := range pairs {
		var krakenPair string
		krakenPair, err = kraken.GetKrakenPair(pair)
		if err != nil {
			return
		}
		krakenPairs = append(krakenPairs, krakenPair)
	}

	// api could not respond, try three times before returning errors
	var count int
	var response *krakenClient.TickerResponse
	response, err = api.client.Ticker(krakenPairs...)
	for err != nil && count < kraken.NbRequestRetries {
		count++
		logger.Errorf("Got error from Kraken api : %s, retry %d/%d after %0.0f seconds", err.Error(), count,
			kraken.NbRequestRetries, kraken.DelayBetweenRetries.Seconds())
		time.Sleep(kraken.DelayBetweenRetries)
		response, err = api.client.Ticker(krakenPairs...)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// search required value from response
	value := reflect.ValueOf(response).Elem()
	typeOfResponse := value.Type()
	for i := 0; i < value.NumField(); i++ {
		pairTickerInfoInterface := value.Field(i).Interface()
		pairTickerInfo, ok := pairTickerInfoInterface.(krakenClient.PairTickerInfo)
		if ok && len(pairTickerInfo.Ask) > 0 && len(pairTickerInfo.Bid) > 0 {
			var pair tradingConst.Pair
			pair, err = kraken.GetProjectPair(typeOfResponse.Field(i).Name)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			var askPrice, bidPrice float64
			askPrice, err = strconv.ParseFloat(pairTickerInfo.Ask[0], 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			bidPrice, err = strconv.ParseFloat(pairTickerInfo.Bid[0], 64)
			if err != nil {
				return nil, errors.WithStack(err)
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
	response, err := api.client.OpenOrders(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for orderID, krakenOrder := range response.Open {
		var order model.Order
		order, err = api.convertOrderFromPlatformToProject(krakenOrder)
		if err != nil {
			return
		}
		order.ID = orderID
		orders = append(orders, order)
	}
	return
}

func (api *krakenApi) GetAllOrders(since time.Time) (orders []model.Order, err error) {
	// retrieve close orders
	response, err := api.client.ClosedOrders(map[string]string{
		kraken.StartCloseOrderParameter: fmt.Sprintf("%d", since.Unix()),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for orderID, krakenOrder := range response.Closed {
		var order model.Order
		order, err = api.convertOrderFromPlatformToProject(krakenOrder)
		if err != nil {
			return
		}
		order.ID = orderID
		orders = append(orders, order)
	}

	// adding open orders
	openOrders, err := api.GetOpenOrders()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orders = append(orders, openOrders...)
	return
}

func (api *krakenApi) RefreshBalance(balance *model.Balance) (err error) {
	// get response form kraken api
	response, err := api.client.Balance()
	if err != nil {
		return errors.WithStack(err)
	}

	// search required value from response
	value := reflect.ValueOf(response).Elem()
	typeOfResponse := value.Type()
	for i := 0; i < value.NumField(); i++ {
		balanceCurrencyInterface := value.Field(i).Interface()
		balanceCurrency, ok := balanceCurrencyInterface.(float64)
		if ok {
			var currency tradingConst.Currency
			currency, ok = kraken.CurrencyConverterFromKrakenToProject[typeOfResponse.Field(i).Name]
			if ok {
				balance.ValueByCurrency[currency] = balanceCurrency
			}
		}
	}

	balance.UpdatedAt = time.Now()
	return
}

func (api *krakenApi) TakerFeesInPercent() (fees float64) {
	return kraken.TakerFees
}

func (api *krakenApi) MakerFeesInPercent() (fees float64) {
	return kraken.MakerFees
}

func (api *krakenApi) convertOrderFromPlatformToProject(krakenOrder krakenClient.Order) (order model.Order, err error) {
	projectPair, err := kraken.GetProjectAssetPair(krakenOrder.Description.AssetPair)
	if err != nil {
		return order, errors.WithStack(err)
	}
	projectSide, err := kraken.GetProjectSide(krakenOrder.Description.Type)
	if err != nil {
		return order, errors.WithStack(err)
	}
	projectStatus, err := kraken.GetProjectStatus(krakenOrder.Status)
	if err != nil {
		return order, errors.WithStack(err)
	}
	projectOrderType, err := kraken.GetProjectOrderType(krakenOrder.Description.OrderType)
	if err != nil {
		return order, errors.WithStack(err)
	}
	var leverage int64
	if krakenOrder.Description.Leverage != "" && krakenOrder.Description.Leverage != "none" {
		leverage, err = strconv.ParseInt(krakenOrder.Description.Leverage, 0, 64)
		if err != nil {
			return order, errors.WithStack(err)
		}
	}
	volume, err := strconv.ParseFloat(krakenOrder.Volume, 64)
	if err != nil {
		return order, errors.WithStack(err)
	}

	var price float64
	if projectStatus == tradingConst.Open {
		price, err = strconv.ParseFloat(krakenOrder.Description.PrimaryPrice, 64)
		if err != nil {
			return
		}
	} else {
		price = krakenOrder.Price
	}

	// find close condition type
	var closeOrderType tradingConst.OrderType
	if strings.Contains(krakenOrder.Description.Close, kraken.CloseTypeDescriptionConverter[tradingConst.Limit]) {
		closeOrderType = tradingConst.Limit
	} else if strings.Contains(krakenOrder.Description.Close, kraken.CloseTypeDescriptionConverter[tradingConst.TakeProfit]) {
		closeOrderType = tradingConst.TakeProfit
	} else if strings.Contains(krakenOrder.Description.Close, kraken.CloseTypeDescriptionConverter[tradingConst.StopLoss]) {
		closeOrderType = tradingConst.StopLoss
	} else {
		closeOrderType = tradingConst.None
	}

	// find close condition price
	var closeConditionPrice float64
	re := regexp.MustCompile("[0-9.]+")
	result := re.FindAllString(krakenOrder.Description.Close, 1)
	if len(result) > 0 {
		closeConditionPrice, err = strconv.ParseFloat(result[0], 64)
		if err != nil {
			return order, errors.WithStack(err)
		}
	}

	order = model.Order{
		OpenTime:            time.Unix(int64(krakenOrder.OpenTime), 0),
		CloseTime:           time.Unix(int64(krakenOrder.CloseTime), 0),
		Pair:                projectPair,
		Side:                projectSide,
		Volume:              volume,
		Type:                projectOrderType,
		Price:               price,
		Amount:              price * volume,
		Leverage:            int(leverage),
		CloseConditionType:  closeOrderType,
		CloseConditionPrice: closeConditionPrice,
		Fees:                krakenOrder.Fee,
		Status:              projectStatus,
		PlatformName:        api.Name(),
	}
	return
}
