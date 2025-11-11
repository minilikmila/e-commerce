package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type ProductFilter struct {
	Search string
	Limit  int
	Offset int
}

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]domain.Product, int64, error)
}
