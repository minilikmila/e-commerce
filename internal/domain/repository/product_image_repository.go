package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type ProductImageRepository interface {
	AddMany(ctx context.Context, images []domain.ProductImage) error
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error)
	CountByProduct(ctx context.Context, productID uuid.UUID) (int64, error)
}
