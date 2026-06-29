package ports

import (
	"context"

	"github.com/FernandoJMartins/microservices/shipping/internal/application/core/domain"
)

type APIPort interface {
	Schedule(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error)
}
