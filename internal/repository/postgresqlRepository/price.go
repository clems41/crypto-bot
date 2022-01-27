package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
)

func (r *repo) StorePrice(price *model.Price) (err error) {
	return
}

func (r *repo) GetPriceHistory(form repository.GetPriceHistoryForm) (prices []model.Price, err error) {
	return
}
