package converter

import (
	"github.com/bogdanovds/rocket_factory/order/internal/model"
	repoModel "github.com/bogdanovds/rocket_factory/order/internal/repository/model"
)

// ToServiceModel конвертирует модель repository в модель сервисного слоя
func ToServiceModel(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}

	return &model.Order{
		ID:            order.ID,
		UserID:        order.UserID,
		PartIDs:       order.PartIDs,
		TotalPrice:    order.TotalPrice,
		Status:        model.OrderStatus(order.Status),
		PaymentMethod: order.PaymentMethod,
		TransactionID: order.TransactionID,
	}
}

// ToRepoModel конвертирует модель сервисного слоя в модель repository
func ToRepoModel(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}

	return &repoModel.Order{
		ID:            order.ID,
		UserID:        order.UserID,
		PartIDs:       order.PartIDs,
		TotalPrice:    order.TotalPrice,
		Status:        repoModel.OrderStatus(order.Status),
		PaymentMethod: order.PaymentMethod,
		TransactionID: order.TransactionID,
	}
}
