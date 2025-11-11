package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	URL       string
	CreatedAt time.Time
}
