package postgresqlRepository

import (
	"crypto-bot/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var _ repository.Repository = (*repo)(nil)

type Metadata struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ExecutionID string
}

type repo struct {
	db          *gorm.DB
	executionID string
}

func New(DB *gorm.DB) (r *repo, err error) {
	r = &repo{
		db:          DB,
		executionID: uuid.New().String(),
	}
	err = r.Migrate()
	if err != nil {
		return
	}
	return
}

func (r *repo) Migrate() (err error) {
	err = r.db.AutoMigrate(&price{}, &balance{}, &order{})
	if err != nil {
		return
	}
	return
}
