package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"github.com/pkg/errors"
)

type config struct {
	Metadata
	Name       string
	Parameters string
}

func (r *repo) StoreConfig(configModel *model.Config) (err error) {
	err = configModel.Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	configRepo := config{
		Name:       configModel.Name,
		Parameters: configModel.Parameters,
	}
	configRepo.ExecutionID = r.executionID

	err = r.db.
		Create(&configRepo).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}
