package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type order struct {
	Metadata
	PlatformOrderID     string `gorm:"uniqueIndex"`
	Pair                string `gorm:"index"`
	Side                string `gorm:"index"`
	Volume              float64
	Type                string
	Price               float64
	Amount              float64
	Leverage            int
	CloseConditionType  string
	CloseConditionPrice float64
	Fees                float64
	Status              string `gorm:"index"`
	PlatformName        string `gorm:"index"`
}

func (r *repo) StoreOrder(orderModel *model.Order) (err error) {
	err = orderModel.Validate()
	if err != nil {
		return
	}

	// find if order already exist
	var orderRepo order
	orderExists := true
	err = r.db.
		Take(&orderRepo, "platform_order_id = ?", orderModel.ID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			orderExists = false
		} else {
			return errors.WithStack(err)
		}
	}

	// get previous order status
	previousStatus := orderRepo.Status

	// fill orderModel fields
	orderRepo.PlatformOrderID = orderModel.ID
	orderRepo.Pair = orderModel.Pair
	orderRepo.Side = orderModel.Side
	orderRepo.Volume = orderModel.Volume
	orderRepo.Type = orderModel.Type
	orderRepo.Price = orderModel.Price
	orderRepo.Amount = orderModel.Amount
	orderRepo.Leverage = orderModel.Leverage
	orderRepo.CloseConditionType = orderModel.CloseConditionType
	orderRepo.CloseConditionPrice = orderModel.CloseConditionPrice
	orderRepo.Fees = orderModel.Fees
	orderRepo.Status = orderModel.Status
	orderRepo.PlatformName = orderModel.PlatformName
	orderRepo.ExecutionID = r.executionID

	// update order if exists, if not create it
	if orderExists {
		// update only if new status
		if previousStatus != orderRepo.Status {
			err = r.db.
				Save(&orderRepo).
				Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
	} else {
		err = r.db.
			Create(&orderRepo).
			Error
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return
}
