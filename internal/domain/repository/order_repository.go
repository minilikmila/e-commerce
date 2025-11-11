package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	HasPendingOrdersByProductID(ctx context.Context, productID uuid.UUID) (bool, error)
}
