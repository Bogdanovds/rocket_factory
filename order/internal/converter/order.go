package converter

import (
	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

func ConvertOrderToDTO(order *model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:  order.ID,
		UserUUID:   order.UserID,
		PartUuids:  order.PartIDs,
		TotalPrice: float32(order.TotalPrice),
		Status:     convertStatusToDTO(order.Status),
		PaymentMethod: orderV1.OptPaymentMethod{
			Value: orderV1.PaymentMethod(order.PaymentMethod),
			Set:   order.PaymentMethod != "",
		},
		TransactionUUID: orderV1.OptNilUUID{
			Value: order.TransactionID,
			Set:   order.TransactionID != uuid.Nil,
			Null:  order.TransactionID == uuid.Nil,
		},
	}
}

func convertStatusToDTO(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPending:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}
