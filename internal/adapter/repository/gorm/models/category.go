package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/minilik/ecommerce/internal/domain"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) ToDomain() *domain.Category {
	// images := make([]domain.ProductImage, 0, len(p.Images))
	// for _, im := range p.Images {
	// 	images = append(images, im.ToDomain())
	// }
	return &domain.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func CategoryFromDomain(cat *domain.Category) *Category {
	if cat == nil {
		return nil
	}
	return &Category{
		ID:        cat.ID,
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt,
		UpdatedAt: cat.UpdatedAt,
	}
}
