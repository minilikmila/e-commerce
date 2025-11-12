package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/adapter/middleware"
	"github.com/minilik/ecommerce/internal/domain"
	orderusecase "github.com/minilik/ecommerce/internal/usecase/order"
)

type mockOrderService struct {
	mock.Mock
}

func (m *mockOrderService) Create(ctx context.Context, userID uuid.UUID, input orderusecase.CreateOrderInput) (*domain.Order, error) {
	args := m.Called(ctx, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *mockOrderService) ListForUser(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Order), args.Error(1)
}

func TestOrderHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockOrderService)
		handler := NewOrderHandler(mockSvc, logger)

		input := orderusecase.CreateOrderInput{
			Items: []orderusecase.OrderItemInput{
				{ProductID: uuid.New(), Quantity: 2},
			},
		}
		order := &domain.Order{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			TotalPrice: 100.0,
			Status:     domain.OrderStatusPending,
		}

		mockSvc.On("Create", mock.Anything, mock.Anything, input).Return(order, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("currentUser", middleware.UserClaims{UserID: uuid.New(), Role: domain.RoleUser})

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestOrderHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockOrderService)
		handler := NewOrderHandler(mockSvc, logger)

		orders := []domain.Order{}

		mockSvc.On("ListForUser", mock.Anything, mock.Anything).Return(orders, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("currentUser", middleware.UserClaims{UserID: uuid.New(), Role: domain.RoleUser})

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

