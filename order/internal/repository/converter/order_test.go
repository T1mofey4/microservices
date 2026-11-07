package converter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/T1mofey4/microservices/order/internal/model"
	repoModel "github.com/T1mofey4/microservices/order/internal/repository/model"
)

func TestStatusToRepoModel(t *testing.T) {
	tests := []struct {
		name     string
		status   model.OrderStatus
		expected repoModel.OrderStatus
	}{
		{name: "PendingPayment", status: model.StatusPendingPayment, expected: repoModel.StatusPendingPayment},
		{name: "Paid", status: model.StatusPaid, expected: repoModel.StatusPaid},
		{name: "Cancelled", status: model.StatusCancelled, expected: repoModel.StatusCancelled},
		{name: "Unknown", status: model.OrderStatus("Unknown"), expected: repoModel.OrderStatus("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := StatusToRepoModel(tt.status)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusFromRepoModel(t *testing.T) {
	tests := []struct {
		name     string
		status   repoModel.OrderStatus
		expected model.OrderStatus
	}{
		{name: "PendingPayment", status: repoModel.StatusPendingPayment, expected: model.StatusPendingPayment},
		{name: "Paid", status: repoModel.StatusPaid, expected: model.StatusPaid},
		{name: "Cancelled", status: repoModel.StatusCancelled, expected: model.StatusCancelled},
		{name: "Unknown", status: repoModel.OrderStatus("Unknown"), expected: model.OrderStatus("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := StatusFromRepoModel(tt.status)
			require.Equal(t, tt.expected, result)
		})
	}
}
