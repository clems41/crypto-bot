package gormUtils

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Metadata gorm.Model

// GetID Return the ID
func (mtd *Metadata) GetID() uint {
	return mtd.ID
}

// GetIDAsString Return the ID in string format
func (mtd *Metadata) GetIDAsString() string {
	return strconv.FormatUint(uint64(mtd.ID), 10)
}

// SetID Set the ID like a string
func (mtd *Metadata) SetID(id uint) {
	mtd.ID = id
}

func (mtd *Metadata) CreatedAtNow() {
	mtd.CreatedAt = time.Now()
}

func (mtd *Metadata) UpdatedAtNow() {
	mtd.CreatedAt = time.Now()
}

func (mtd *Metadata) HasBeenModified() bool {

	if mtd == nil {
		return false
	}

	return !mtd.UpdatedAt.IsZero() && !mtd.UpdatedAt.Equal(mtd.CreatedAt)
}

// Pagination is intended to be embedded in objects targeted by paginated queries
// In this case, GetCountOver is intended to be used in the ''Select' SQL clause
// See
type Pagination struct {
	TotalPagination int `gorm:"-"`
}

// GetCountOver Returns a gorm SQL expression that counts total records into TotalPagination field
// This is useful for request with pagination.
// Usage DB.Select("?, *", encounter.GetCountOver())
func (pagination Pagination) GetCountOver() string {
	return "count(*) OVER() AS total_pagination"
}
