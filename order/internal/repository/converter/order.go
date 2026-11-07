package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
	repoModel "github.com/T1mofey4/microservices/order/internal/repository/model"
)

func StatusToRepoModel(s model.OrderStatus) repoModel.OrderStatus {
	switch s {
	case model.StatusPendingPayment:
		return repoModel.StatusPendingPayment
	case model.StatusPaid:
		return repoModel.StatusPaid
	case model.StatusCancelled:
		return repoModel.StatusCancelled
	default:
		return repoModel.OrderStatus("")
	}
}

func StatusFromRepoModel(s repoModel.OrderStatus) model.OrderStatus {
	switch s {
	case repoModel.StatusPendingPayment:
		return model.StatusPendingPayment
	case repoModel.StatusPaid:
		return model.StatusPaid
	case repoModel.StatusCancelled:
		return model.StatusCancelled
	default:
		return model.OrderStatus("")
	}
}

func PaymentMethodToRepoModel(p *model.PaymentMethod) *repoModel.PaymentMethod {
	if p == nil {
		return nil
	}

	var method repoModel.PaymentMethod

	switch *p {
	case model.PaymentMethodUnknown:
		method = repoModel.PaymentMethodUnknown
	case model.PaymentMethodCard:
		method = repoModel.PaymentMethodCard
	case model.PaymentMethodSBP:
		method = repoModel.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		method = repoModel.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		method = repoModel.PaymentMethodInvestorMoney
	default:
		method = repoModel.PaymentMethodUnknown
	}

	return &method
}

func PaymentMethodFromRepoModel(p *repoModel.PaymentMethod) *model.PaymentMethod {
	if p == nil {
		return nil
	}

	var method model.PaymentMethod

	switch *p {
	case repoModel.PaymentMethodUnknown:
		method = model.PaymentMethodUnknown
	case repoModel.PaymentMethodCard:
		method = model.PaymentMethodCard
	case repoModel.PaymentMethodSBP:
		method = model.PaymentMethodSBP
	case repoModel.PaymentMethodCreditCard:
		method = model.PaymentMethodCreditCard
	case repoModel.PaymentMethodInvestorMoney:
		method = model.PaymentMethodInvestorMoney
	default:
		method = model.PaymentMethodUnknown
	}

	return &method
}

func OrderToModel(o repoModel.Order) (model.Order, error) {
	orderUUID, err := uuid.Parse(o.OrderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("incorrect order uuid: %w", err)
	}

	userUUID, err := uuid.Parse(o.UserUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("incorrect user uuid: %w", err)
	}

	partUUIDs := make([]uuid.UUID, 0, len(o.PartUUIDs))
	for _, p := range o.PartUUIDs {
		partUUID, err := uuid.Parse(p)
		if err != nil {
			return model.Order{}, fmt.Errorf("incorrect part uuid: %w", err)
		}
		partUUIDs = append(partUUIDs, partUUID)
	}

	var transactionUUID *uuid.UUID
	if o.TransactionUUID != nil {
		t, err := uuid.Parse(*o.TransactionUUID)
		if err != nil {
			return model.Order{}, fmt.Errorf("incorrect transaction uuid: %w", err)
		}

		transactionUUID = &t
	}

	return model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   PaymentMethodFromRepoModel(o.PaymentMethod),
		Status:          StatusFromRepoModel(o.Status),
	}, nil
}

func OrderToRepoModel(o model.Order) repoModel.Order {
	partUUIDs := make([]string, 0, len(o.PartUUIDs))
	for _, partUUID := range o.PartUUIDs {
		partUUIDs = append(partUUIDs, partUUID.String())
	}

	var transactionUUID *string
	if o.TransactionUUID != nil {
		t := o.TransactionUUID.String()
		transactionUUID = &t
	}

	return repoModel.Order{
		OrderUUID:       o.OrderUUID.String(),
		UserUUID:        o.UserUUID.String(),
		PartUUIDs:       partUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   PaymentMethodToRepoModel(o.PaymentMethod),
		Status:          StatusToRepoModel(o.Status),
	}
}
