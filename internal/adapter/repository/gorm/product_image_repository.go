package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/minilik/ecommerce/internal/adapter/repository/gorm/models"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
)

type productImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) repository.ProductImageRepository {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) AddMany(ctx context.Context, images []domain.ProductImage) error {
	if len(images) == 0 {
		return nil
	}
	rows := make([]models.ProductImage, 0, len(images))
	now := time.Now()
	for _, img := range images {
		id := img.ID
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.ProductImage{
			ID:        id,
			ProductID: img.ProductID,
			URL:       img.URL,
			CreatedAt: now,
		})
	}
	return r.db.WithContext(ctx).Create(&rows).Error
}

func (r *productImageRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error) {
	var rows []models.ProductImage
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ProductImage, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.ToDomain())
	}
	return out, nil
}

func (r *productImageRepository) CountByProduct(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
