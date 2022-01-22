package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/pkg/utils/csvUtils"
	"crypto-bot/pkg/utils/pathUtils"
	"fmt"
	"time"
)

var _ repository.Balance = (*balanceRepo)(nil)

type balanceRepo struct {
	balances    map[string]map[string]*repositoryModel.Balance // store all balance by platform and currency
	csvFilePath string
}

func NewBalanceRepository() (repo *balanceRepo, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}
	csvPath := fmt.Sprintf("%s/%s", rootPath, csvFileNameBalance)

	// empty file or creating new one
	err = csvUtils.ForceCreateFile(csvPath)
	if err != nil {
		return
	}

	// adding column names
	columnNames := []string{
		"PlatformName",
		"Currency",
		"Value",
		"UpdatedAt",
	}
	err = csvUtils.AppendLines(csvPath, columnNames)
	if err != nil {
		return
	}

	repo = &balanceRepo{
		balances:    make(map[string]map[string]*repositoryModel.Balance),
		csvFilePath: csvPath,
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

	// update csv line
	newLine := []string{
		balance.PlatformName,
		balance.Currency,
		fmt.Sprintf("%0.2f", balance.Value),
		balance.UpdatedAt.Format(time.RFC3339),
	}

	// update line or append it if not exists
	matchingString := fmt.Sprintf("%s,%s", balance.PlatformName, balance.Currency)
	err = csvUtils.UpdateLineOrAppend(repo.csvFilePath, matchingString, newLine)
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
