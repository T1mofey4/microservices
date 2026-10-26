package converter

import (
	"fmt"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func OrderToOrderDto(order model.Order) orderV1.OrderDto {
	var transactionUUID orderV1.OptNilUUID
	if order.TransactionUUID == nil {
		transactionUUID = orderV1.OptNilUUID{Set: true, Null: true}
	} else {
		transactionUUID = orderV1.NewOptNilUUID(*order.TransactionUUID)
	}

	return orderV1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   PaymentMethodToOrderDTO(order.PaymentMethod),
		Status:          OrderStatusToOrderDTO(order.Status),
	}
}

func OrderStatusToOrderDTO(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.StatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.StatusPaid:
		return orderV1.OrderStatusPAID
	case model.StatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatus("")
	}
}

func PaymentMethodToOrderDTO(method *model.PaymentMethod) orderV1.OptNilPaymentMethod {
	if method == nil {
		return orderV1.OptNilPaymentMethod{Set: true, Null: true}
	}

	switch *method {
	case model.PaymentMethodUnknown:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodUNKNOWN)
	case model.PaymentMethodCard:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodCARD)
	case model.PaymentMethodSBP:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodSBP)
	case model.PaymentMethodCreditCard:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodCREDITCARD)
	case model.PaymentMethodInvestorMoney:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodINVESTORMONEY)
	default:
		return orderV1.NewOptNilPaymentMethod(orderV1.PaymentMethodUNKNOWN)
	}
}

func OrderToCreateOrderResponse(order model.Order) orderV1.CreateOrderResponse {
	return orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
}

func PaymentMethodToModel(method orderV1.PaymentMethod) (model.PaymentMethod, error) {
	switch method {
	case orderV1.PaymentMethodUNKNOWN:
		return model.PaymentMethod(""), fmt.Errorf("unknown payment method: %w", model.ErrInvalidPaymentMethod)
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethod(""), fmt.Errorf("unknown payment method: %w", model.ErrInvalidPaymentMethod)
	}
}
