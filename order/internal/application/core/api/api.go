package api

import (
	"github.com/FernandoJMartins/microservices/order/internal/application/core/domain"
	"github.com/FernandoJMartins/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{db: db, payment: payment, shipping: shipping}
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

	for _, item := range order.OrderItems {
		product, err := a.db.GetProductByCode(item.ProductCode)
		if err != nil {
			return domain.Order{}, status.Errorf(codes.NotFound, "Product %s not found.", item.ProductCode)
		}
		if product.Quantity < item.Quantity {
			return domain.Order{}, status.Errorf(codes.InvalidArgument, "Product %s has insufficient stock. Available: %d, Requested: %d.", item.ProductCode, product.Quantity, item.Quantity)
		}
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

	shippingErr := a.shipping.Schedule(&order)
	if shippingErr != nil {
		order.Status = "Cancelled"
		a.db.Save(&order)
		return domain.Order{}, shippingErr
	}

	return order, nil
}
