package googleSheetRepository

import (
	"context"
	"crypto-bot/internal/repository"
	"google.golang.org/api/sheets/v4"
	"time"
)

var _ repository.Repository = (*repo)(nil)

type repo struct {
	googleSheetService *sheets.Service
}

func New(ctx context.Context) (r *repo, err error) {
	googleSheetService, err := sheets.NewService(ctx)
	if err != nil {
		return
	}
	r = &repo{
		googleSheetService: googleSheetService,
	}
	return
}

func (r *repo) StorePrice(price *model.Price) (err error) {
	err = price.Validate()
	if err != nil {
		return
	}

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				price.Date.Format(time.RFC3339),
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

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				balance.PlatformName,
				balance.ValueByCurrency,
				balance.UpdatedAt.Format(time.RFC3339),
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
