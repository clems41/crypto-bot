package repositoryModel

import "time"

type Balance struct {
	PlatformName string
	Currency     string
	Value        float64
	UpdatedAt    time.Time
}
