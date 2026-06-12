package ports

import "github.com/FernandoJMartins/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(order domain.Order) error
	// Get(id string) (domain.Order, error)
	// Save(*domain.Order) error
}
