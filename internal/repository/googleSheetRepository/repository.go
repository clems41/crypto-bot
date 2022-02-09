package googleSheetRepository

import (
	"context"
	"crypto-bot/internal/constant/timeConst"
	"crypto-bot/internal/constant/tradingConst"
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"reflect"
)

var _ repository.Repository = (*repo)(nil)

type repo struct {
	googleSheetService     *sheets.Service
	previousBalance        map[tradingConst.Currency]float64
	pricesByPlatformByPair map[string]map[tradingConst.Pair][]model.Price
}

func New(ctx context.Context) (r *repo, err error) {
	googleSheetService, err := sheets.NewService(ctx)
	if err != nil {
		return
	}
	r = &repo{
		googleSheetService:     googleSheetService,
		previousBalance:        make(map[tradingConst.Currency]float64),
		pricesByPlatformByPair: make(map[string]map[tradingConst.Pair][]model.Price),
	}
	return
}

func (r *repo) StorePrice(price *model.Price) (err error) {
	err = price.Validate()
	if err != nil {
		return
	}

	// store in memory
	if r.pricesByPlatformByPair[price.PlatformName] == nil {
		r.pricesByPlatformByPair[price.PlatformName] = make(map[tradingConst.Pair][]model.Price)
	}
	r.pricesByPlatformByPair[price.PlatformName][price.Pair] = append(
		r.pricesByPlatformByPair[price.PlatformName][price.Pair], *price)

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				price.Date.Format(timeConst.DefaultFormatTimeLayout),
				price.PlatformName,
				price.Pair,
				price.Ask,
				price.Bid,
			},
		},
	}
	_, err = r.googleSheetService.Spreadsheets.Values.
		Append(googleSpreadsheetID, priceRangeUpdate, &valueRange).
		ValueInputOption("RAW").
		Do()
	return
}

func (r *repo) StoreBalance(balance *model.Balance) (err error) {
	err = balance.Validate()
	if err != nil {
		return
	}

	// don't send new balance into google spreadsheet if no new changes
	if !reflect.DeepEqual(balance.ValueByCurrency, r.previousBalance) {
		// append new line in google spreadsheet
		valueRange := sheets.ValueRange{
			Values: [][]interface{}{
				{
					balance.PlatformName,
					fmt.Sprintf("%v", balance.ValueByCurrency),
					balance.UpdatedAt.Format(timeConst.DefaultFormatTimeLayout),
				},
			},
		}
		_, err = r.googleSheetService.Spreadsheets.Values.
			Append(googleSpreadsheetID, balanceRangeUpdate, &valueRange).
			ValueInputOption("RAW").
			Do()
		if err != nil {
			return
		}
	}

	// do this instead of directly copy to avoid getting pointer
	for currency, value := range balance.ValueByCurrency {
		r.previousBalance[currency] = value
	}

	return
}

func (r *repo) StoreOrder(order *model.Order) (err error) {
	err = order.Validate()
	if err != nil {
		return
	}

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				order.ID,
				order.OpenTime,
				order.CloseTime,
				order.Pair,
				order.Side,
				order.Volume,
				order.Type,
				order.Price,
				order.Amount,
				order.Leverage,
				order.CloseConditionType,
				order.CloseConditionPrice,
				order.Fees,
				order.Status,
			},
		},
	}
	_, err = r.googleSheetService.Spreadsheets.Values.
		Append(googleSpreadsheetID, orderRangeUpdate, &valueRange).
		ValueInputOption("RAW").
		Do()
	return
}

func (r *repo) GetPriceHistory(form repository.GetPriceHistoryForm) (prices []model.Price, err error) {
	priceHistoryByPlatform, ok := r.pricesByPlatformByPair[form.PlatformName]
	if !ok {
		return prices, fmt.Errorf("cannot find price history for platform %s", form.PlatformName)
	}
	priceHistory, ok := priceHistoryByPlatform[form.Pair]
	if !ok {
		return prices, fmt.Errorf("cannot find price history for pair %s", form.Pair)
	}
	for _, price := range priceHistory {
		if price.Date.After(form.SinceTime) {
			prices = append(prices, price)
		}
	}
	return
}

func (r *repo) GetOrderHistory(form repository.GetOrderHistoryForm) (orders []model.Order, err error) {
	// TODO
	return
}

func (r *repo) StoreConfig(configModel *model.Config) (err error) {
	// TODO
	return
}