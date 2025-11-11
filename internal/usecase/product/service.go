package product

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
	Create(ctx context.Context, ownerID uuid.UUID, input CreateProductInput) (*domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateProductInput) (*domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, input ListProductsInput) ([]domain.Product, int64, error)
}

type service struct {
	repo      repository.ProductRepository
	orderRepo repository.OrderRepository
	logger    *zap.Logger
	now       func() time.Time
}

func NewService(repo repository.ProductRepository, orderRepo repository.OrderRepository, logger *zap.Logger) Service {
	return &service{
		repo:      repo,
		orderRepo: orderRepo,
		logger:    logger,
		now:       time.Now,
	}
}

func (s *service) Create(ctx context.Context, ownerID uuid.UUID, input CreateProductInput) (*domain.Product, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	product := &domain.Product{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    strings.TrimSpace(input.Category),
		UserID:      ownerID,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateProductInput) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	if err := applyUpdate(product, input); err != nil {
		return nil, err
	}

	product.UpdatedAt = s.now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if product exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrProductNotFound
	}

	// Check if there are any pending orders for this product
	hasPending, err := s.orderRepo.HasPendingOrdersByProductID(ctx, id)
	if err != nil {
		s.logger.Error("failed to check pending orders for product", zap.String("product_id", id.String()), zap.Error(err))
		return fmt.Errorf("failed to check pending orders: %w", err)
	}

	if hasPending {
		return domain.ErrProductHasPendingOrders
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (s *service) List(ctx context.Context, input ListProductsInput) ([]domain.Product, int64, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	filter := repository.ProductFilter{
		Search: strings.TrimSpace(input.Search),
		Limit:  pageSize,
		Offset: offset,
	}

	return s.repo.List(ctx, filter)
}

func validateCreateInput(input CreateProductInput) error {
	if len(strings.TrimSpace(input.Name)) < 3 || len(strings.TrimSpace(input.Name)) > 100 {
		return fmt.Errorf("required:name must be between 3 and 100 characters")
	}
	if len(strings.TrimSpace(input.Description)) < 10 {
		return fmt.Errorf("required:description must be at least 10 characters")
	}
	if input.Price <= 0 {
		return fmt.Errorf("required:price must be greater than zero")
	}
	if input.Stock < 0 {
		return fmt.Errorf("required:stock must be non-negative")
	}
	if strings.TrimSpace(input.Category) == "" {
		return fmt.Errorf("required:category is required")
	}
	return nil
}

func applyUpdate(product *domain.Product, input UpdateProductInput) error {
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if len(name) == 0 {
			return fmt.Errorf("name cannot be empty")
		}
		product.Name = name
	}
	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		if len(desc) == 0 {
			return fmt.Errorf("description cannot be empty")
		}
		product.Description = desc
	}
	if input.Price != nil {
		if *input.Price <= 0 {
			return fmt.Errorf("price must be greater than zero")
		}
		product.Price = *input.Price
	}
	if input.Stock != nil {
		if *input.Stock < 0 {
			return fmt.Errorf("stock must be non-negative")
		}
		product.Stock = *input.Stock
	}
	if input.Category != nil {
		category := strings.TrimSpace(*input.Category)
		if category == "" {
			return fmt.Errorf("category cannot be empty")
		}
		product.Category = category
	}
	return nil
}
