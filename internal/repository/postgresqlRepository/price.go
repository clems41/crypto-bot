package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
	"github.com/pkg/errors"
)

type price struct {
	Metadata
	PlatformName string  `gorm:"index"`
	Pair         string  `gorm:"index"`
	Ask          float64 `gorm:"index"`
	Bid          float64 `gorm:"index"`
}

func (r *repo) StorePrice(priceModel *model.Price) (err error) {
	err = priceModel.Validate()
	if err != nil {
		return
	}

	priceRepo := price{
		PlatformName: priceModel.PlatformName,
		Pair:         priceModel.Pair,
		Ask:          priceModel.Ask,
		Bid:          priceModel.Bid,
	}
	priceRepo.ExecutionID = r.executionID

	err = r.db.
		Create(&priceRepo).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return
}

func (r *repo) GetPriceHistory(form repository.GetPriceHistoryForm) (prices []model.Price, err error) {
	var pricesRepo []price
	err = r.db.
		Find(&pricesRepo, "execution_id = ? AND platform_name = ? AND pair = ? AND created_at >= ?",
			r.executionID, form.PlatformName, form.Pair, form.SinceTime).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// fill prices model
	for _, priceRepo := range pricesRepo {
		priceModel := model.Price{
			Date:         priceRepo.CreatedAt,
			PlatformName: priceRepo.PlatformName,
			Pair:         priceRepo.Pair,
			Ask:          priceRepo.Ask,
			Bid:          priceRepo.Bid,
		}
		prices = append(prices, priceModel)
	}
	return
}
