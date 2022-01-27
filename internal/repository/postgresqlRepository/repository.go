package postgresqlRepository

import (
	"crypto-bot/internal/repository"
	"gorm.io/gorm"
)

var _ repository.Repository = (*repo)(nil)

type repo struct {
	db *gorm.DB
}

func New(DB *gorm.DB) (r *repo, err error) {
	r = &repo{
		db: DB,
	}
	return
}
