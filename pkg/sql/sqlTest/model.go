package sqlTest

import "gorm.io/gorm"

type Repositories interface {
	Migrate(DB *gorm.DB) (err error)
}
