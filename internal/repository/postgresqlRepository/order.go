package postgresqlRepository

import (
	"crypto-bot/internal/model"
	"crypto-bot/internal/repository"
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
		return errors.WithStack(err)
	}

	// find if order already exist
	var orderRepo order
	orderExists := true
	err = r.db.
		Take(&orderRepo, "execution_id = ? AND platform_order_id = ?", r.executionID, orderModel.ID).
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

func (r *repo) GetOrderHistory(form repository.GetOrderHistoryForm) (orders []model.Order, err error) {
	var ordersRepo []order
	query := r.db.Where("execution_id = ?", r.executionID)
	if form.Status != "" {
		query = query.Where("status = ?", form.Status)
	}
	if form.Pair != "" {
		query = query.Where("pair = ?", form.Pair)
	}
	if form.PlatformName != "" {
		query = query.Where("platform_name = ?", form.PlatformName)
	}
	err = query.
		Find(&ordersRepo).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, orderRepo := range ordersRepo {
		orders = append(orders, model.Order{
			ID:                  orderRepo.PlatformOrderID,
			Date:                orderRepo.CreatedAt,
			Pair:                orderRepo.Pair,
			Side:                orderRepo.Side,
			Volume:              orderRepo.Volume,
			Type:                orderRepo.Type,
			Price:               orderRepo.Price,
			Amount:              orderRepo.Amount,
			Leverage:            orderRepo.Leverage,
			CloseConditionType:  orderRepo.CloseConditionType,
			CloseConditionPrice: orderRepo.CloseConditionPrice,
			Fees:                orderRepo.Fees,
			Status:              orderRepo.Status,
			PlatformName:        orderRepo.PlatformName,
		})
	}

	return
}
