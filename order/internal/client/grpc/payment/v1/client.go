package v1

import (
	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"

	def "github.com/T1mofey4/microservices/order/internal/client/grpc"
)

var _ def.PaymentClient = (*paymentClient)(nil)

type paymentClient struct {
	client paymentV1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *paymentClient {
	return &paymentClient{
		client: paymentV1.NewPaymentServiceClient(conn),
	}
}
