package payment_adapter

import (
	"context"
	"log"
	"time"

	"github.com/FernandoJMartins/microservices-proto/golang/payment"
	"github.com/FernandoJMartins/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	payment payment.PaymentClient // é gerado automaticamente pelo compilador do protobuf
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(paymentServiceUrl, opts...)

	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn) // inicia o stub do cliente gRPC
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // garante que o contexto seja cancelado após a operação, evitando vazamentos de recursos

	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	if status.Code(err) == codes.DeadlineExceeded {
		log.Printf("Timeout ao tentar processar o pagamento para o pedido %s", order.ID)
		return nil
	}

	return err
}
