package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;index;not null"`
	URL       string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}

func (ProductImage) TableName() string {
	return "product_images"
}

func (m *ProductImage) ToDomain() domain.ProductImage {
	return domain.ProductImage{
		ID:        m.ID,
		ProductID: m.ProductID,
		URL:       m.URL,
		CreatedAt: m.CreatedAt,
	}
}
