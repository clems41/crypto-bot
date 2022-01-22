package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"google.golang.org/api/sheets/v4"
	"time"
)

var _ repository.Balance = (*balanceRepo)(nil)

type balanceRepo struct {
	balances           map[string]map[string]*repositoryModel.Balance // store all balance by platform and currency
	googleSheetService *sheets.Service
}

func NewBalanceRepository(googleSheetService *sheets.Service) (repo *balanceRepo, err error) {
	repo = &balanceRepo{
		balances:           make(map[string]map[string]*repositoryModel.Balance),
		googleSheetService: googleSheetService,
	}
	return
}

func (repo *balanceRepo) Update(balance *repositoryModel.Balance) (err error) {
	if balance != nil {
		balance.UpdatedAt = time.Now()
		if repo.balances[balance.PlatformName] == nil {
			repo.balances[balance.PlatformName] = make(map[string]*repositoryModel.Balance)
		}
		repo.balances[balance.PlatformName][balance.Currency] = balance
	}

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				balance.PlatformName,
				balance.Currency,
				balance.Value,
				balance.UpdatedAt.Format(time.RFC3339),
			},
		},
	}
	request := repo.googleSheetService.Spreadsheets.Values.
		Append(googleSpreadsheetID, balanceRangeUpdate, &valueRange).
		ValueInputOption("RAW")
	_, err = request.Do()
	if err != nil {
		return
	}
	return
}

func (repo *balanceRepo) GetCurrentBalance(platformName string, currency string) (balance *repositoryModel.Balance, err error) {
	balance, ok := repo.balances[platformName][currency]
	if !ok {
		return nil, repository.ErrBalanceNotFound
	}
	return
}
