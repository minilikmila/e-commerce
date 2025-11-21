package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateCategory) (*domain.Category, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateCategoryInput) (*domain.Category, error)
	List(ctx context.Context, input ListCategoryInput) ([]domain.Category, error)
}

type service struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
	now          func() time.Time
}

func NewService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, logger *zap.Logger) Service {
	return &service{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
		now:          time.Now,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, input CreateCategory) (*domain.Category, error) {
	return nil, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateCategoryInput) (*domain.Category, error) {
	return nil, nil
}

func (s *service) List(ctx context.Context, input ListCategoryInput) ([]domain.Category, error) {
	return nil, nil
}
