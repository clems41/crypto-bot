package postgresqlRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/pkg/logger"
	"crypto-bot/pkg/utils/envUtils"
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

	// cleanup table if env variable is set
	cleanup := envUtils.GetFromEnvOrDefault(envCleanupRepo, defaultCleanupRepo)
	if cleanup == "true" {
		err = r.cleanupTables()
		if err != nil {
			return
		}
	}

	logger.Infof("Execution ID %s", r.executionID)

	return
}

func (r *repo) Migrate() (err error) {
	err = r.db.AutoMigrate(&price{}, &balance{}, &order{})
	if err != nil {
		return
	}
	return
}
