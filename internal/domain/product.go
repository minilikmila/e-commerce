package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity.
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string
	UserID      uuid.UUID
	Images      []ProductImage `json:"images,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
