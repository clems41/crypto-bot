package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type balance struct {
	Metadata
	PlatformName string  `gorm:"index"`
	Value        float64 `gorm:"index"`
	Currency     string  `gorm:"index"`
}

func (r *repo) StoreBalance(balanceModel *model.Balance) (err error) {
	err = balanceModel.Validate()
	if err != nil {
		return
	}

	for currency, value := range balanceModel.ValueByCurrency {
		// if previous value is the same, don't do anything
		var balanceRepo balance
		err = r.db.
			Take(&balanceRepo, "execution_id = ? AND platform_name = ? AND currency = ? AND value = ?",
				r.executionID, balanceModel.PlatformName, currency, value).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// insert new balance value if same previous one has not been found
				balanceRepo.PlatformName = balanceModel.PlatformName
				balanceRepo.Currency = currency
				balanceRepo.Value = value
				balanceRepo.ExecutionID = r.executionID
				err = r.db.
					Create(&balanceRepo).
					Error
				if err != nil {
					return errors.WithStack(err)
				}
			} else {
				return errors.WithStack(err)
			}
		}
	}

	return
}
