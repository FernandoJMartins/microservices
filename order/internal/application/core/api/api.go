package api

import (
	"github.com/FernandoJMartins/microservices/order/internal/application/core/domain"
	"github.com/FernandoJMartins/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{db: db, payment: payment}
}

func (a *Application) PlaceOrder(order domain.Order) (domain.Order, error) {

	var totalItems int32

	for _, item := range order.OrderItems {
		totalItems += item.Quantity
	}

	if totalItems > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Order with more than 50 items is not allowed.")
	}

	if totalItems == 0 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Order cannot have zero items.")
	}

	err := a.db.Save(&order)

	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)

	if paymentErr != nil {
		order.Status = "Cancelled"
	} else {
		order.Status = "Paid"
	}

	err = a.db.Save(&order)

	if err != nil {
		return domain.Order{}, err
	}

	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	return order, nil
}
