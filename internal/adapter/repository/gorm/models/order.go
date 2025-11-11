package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type Order struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Description string    `gorm:"type:text"`
	TotalPrice  float64   `gorm:"not null"`
	Status      string    `gorm:"size:50;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	UnitPrice float64   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (o *Order) ToDomain() *domain.Order {
	items := make([]domain.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, domain.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			OrderID:   item.OrderID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &domain.Order{
		ID:          o.ID,
		UserID:      o.UserID,
		Description: o.Description,
		TotalPrice:  o.TotalPrice,
		Status:      domain.OrderStatus(o.Status),
		Items:       items,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func OrderFromDomain(order *domain.Order) *Order {
	if order == nil {
		return nil
	}

	items := make([]OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderItem{
			ID:        item.ID,
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &Order{
		ID:          order.ID,
		UserID:      order.UserID,
		Description: order.Description,
		TotalPrice:  order.TotalPrice,
		Status:      string(order.Status),
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
