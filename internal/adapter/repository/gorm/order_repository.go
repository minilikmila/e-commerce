package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/minilik/ecommerce/internal/adapter/repository/gorm/models"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	model := models.OrderFromDomain(order)
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	order.ID = model.ID
	return nil
}

func (r *orderRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	var records []models.Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	// it already under session based execution, so no need to create a new transaction
	orders := make([]domain.Order, 0, len(records))
	for _, rec := range records {
		if o := rec.ToDomain(); o != nil {
			orders = append(orders, *o)
		}
	}
	return orders, nil
}

func (r *orderRepository) HasPendingOrdersByProductID(ctx context.Context, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.OrderItem{}).
		Joins("INNER JOIN orders ON order_items.order_id = orders.id").
		Where("order_items.product_id = ? AND orders.status = ?", productID, string(domain.OrderStatusPending)).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
