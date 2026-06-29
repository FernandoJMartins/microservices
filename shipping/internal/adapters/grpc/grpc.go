package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/FernandoJMartins/microservices-proto/golang/shipping"
	"github.com/FernandoJMartins/microservices/shipping/internal/application/core/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	log.Println("Creating shipping...")

	var items []domain.ShippingItem
	for _, item := range request.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	newShipping := domain.NewShipping(request.OrderId, items)
	result, err := a.api.Schedule(ctx, newShipping)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to create shipping. %v", err)).Err()
	}

	return &shipping.CreateShippingResponse{
		ShippingId:   result.ID,
		DeliveryDays: result.DeliveryDays,
	}, nil
}
