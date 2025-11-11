package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateOrderInput) (*domain.Order, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
}

type service struct {
	uow    repository.UnitOfWork
	logger *zap.Logger
	now    func() time.Time
}

func NewService(uow repository.UnitOfWork, logger *zap.Logger) Service {
	return &service{
		uow:    uow,
		logger: logger,
		now:    time.Now,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, input CreateOrderInput) (*domain.Order, error) {
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	order := &domain.Order{
		ID:          uuid.New(),
		UserID:      userID,
		Description: strings.TrimSpace(input.Description),
		Status:      domain.OrderStatusPending,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}
	// Session based transaction
	// This is more efficient than using a single transaction for the entire order creation
	// because it allows for more granular control over the transaction boundaries

	err := s.uow.Execute(ctx, func(repos repository.RepositoryProvider) error {
		var total float64
		items := make([]domain.OrderItem, 0, len(input.Items))

		for _, item := range input.Items {
			if item.Quantity <= 0 {
				return fmt.Errorf("quantity for product %s must be greater than zero", item.ProductID)
			}

			product, err := repos.Products().GetByID(ctx, item.ProductID)
			if err != nil {
				return domain.ErrProductNotFound
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("%w: %s", domain.ErrInsufficientStock, product.Name)
			}

			product.Stock -= item.Quantity
			product.UpdatedAt = s.now()

			if err := repos.Products().Update(ctx, product); err != nil {
				return err
			}

			itemTotal := product.Price * float64(item.Quantity)
			total += itemTotal

			items = append(items, domain.OrderItem{
				ID:        uuid.New(),
				ProductID: product.ID,
				OrderID:   order.ID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
				CreatedAt: s.now(),
				UpdatedAt: s.now(),
			})
		}

		order.TotalPrice = total
		order.Items = items

		if err := repos.Orders().Create(ctx, order); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order
	err := s.uow.Execute(ctx, func(repos repository.RepositoryProvider) error {
		var err error
		orders, err = repos.Orders().ListByUser(ctx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}
