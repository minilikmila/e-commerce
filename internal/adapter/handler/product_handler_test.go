package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/domain"
	productusecase "github.com/minilik/ecommerce/internal/usecase/product"
)

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) Create(ctx context.Context, ownerID uuid.UUID, input productusecase.CreateProductInput) (*domain.Product, error) {
	args := m.Called(ctx, ownerID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductService) Update(ctx context.Context, id uuid.UUID, input productusecase.UpdateProductInput) (*domain.Product, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *mockProductService) List(ctx context.Context, input productusecase.ListProductsInput) ([]domain.Product, int64, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func TestProductHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockProductService)
		handler := NewProductHandler(mockSvc, logger)

		input := productusecase.ListProductsInput{Page: 1, PageSize: 10}
		products := []domain.Product{}
		total := int64(0)

		mockSvc.On("List", mock.Anything, input).Return(products, total, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products?page=1&limit=10", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

