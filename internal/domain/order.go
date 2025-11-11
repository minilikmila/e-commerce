package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents a single line item inside an order.
type OrderItem struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	OrderID   uuid.UUID
	Quantity  int
	UnitPrice float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Order represents an order entity.
type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description string
	TotalPrice  float64
	Status      OrderStatus
	Items       []OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
