package repository

import (
	"context"

	"github.com/minilik/ecommerce/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	List(ctx context.Context, filter ProductFilter) ([]domain.Category, int64, error)
}
