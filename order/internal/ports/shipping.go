package ports

import "github.com/FernandoJMartins/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Schedule(order *domain.Order) error
}
