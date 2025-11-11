package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text;not null"`
	Price       float64   `gorm:"not null"`
	Stock       int       `gorm:"not null"`
	Category    string    `gorm:"size:100;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) ToDomain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		UserID:      p.UserID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ProductFromDomain(product *domain.Product) *Product {
	if product == nil {
		return nil
	}
	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
