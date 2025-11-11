package gorm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/minilik/ecommerce/internal/adapter/repository/gorm/models"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	model := models.ProductFromDomain(product)
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	product.ID = model.ID
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	data := map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
		"category":    product.Category,
		"user_id":     product.UserID,
		"updated_at":  product.UpdatedAt,
	}
	result := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var model models.Product
	if err := r.db.WithContext(ctx).Preload("Images").First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *productRepository) List(ctx context.Context, filter repository.ProductFilter) ([]domain.Product, int64, error) {
	var (
		productList []models.Product
		total       int64
	)

	tx := r.db.WithContext(ctx).Model(&models.Product{})
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		tx = tx.Where("LOWER(name) LIKE ?", search)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		tx = tx.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		tx = tx.Offset(filter.Offset)
	}

	if err := tx.Preload("Images").Order("created_at DESC").Find(&productList).Error; err != nil {
		return nil, 0, err
	}
	// it already under session based execution, so no need to create a new transaction
	// This will be optimized to do more efficient mapping later if needed
	products := make([]domain.Product, 0, len(productList))
	for _, model := range productList {
		if domainProduct := model.ToDomain(); domainProduct != nil {
			products = append(products, *domainProduct)
		}
	}

	return products, total, nil
}
