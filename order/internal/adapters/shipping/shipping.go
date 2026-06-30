package shipping_adapter

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"github.com/FernandoJMartins/microservices-proto/golang/shipping"
	"github.com/FernandoJMartins/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient // é gerado automaticamente pelo compilador do protobuf
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
	)))

	// opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(shippingServiceUrl, opts...)

	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn) // inicia o stub do cliente gRPC
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Schedule(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var items []*shipping.ShippingItem
	for _, item := range order.OrderItems {
		items = append(items, &shipping.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	_, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items:   items,
	})

	if status.Code(err) == codes.DeadlineExceeded {
		log.Printf("Timeout ao tentar processar o envio para o pedido %d", order.ID)
		return err
	}

	return err
}
