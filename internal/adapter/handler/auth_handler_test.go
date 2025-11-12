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

	authusecase "github.com/minilik/ecommerce/internal/usecase/auth"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, input authusecase.RegisterInput) (*authusecase.RegisterResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authusecase.RegisterResponse), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, input authusecase.LoginInput) (*authusecase.AuthResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authusecase.AuthResponse), args.Error(1)
}

func (m *mockAuthService) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockAuthService)
		handler := NewAuthHandler(mockSvc, logger)

		input := authusecase.RegisterInput{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Test123!@#",
		}
		resp := &authusecase.RegisterResponse{
			UserID:   uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		mockSvc.On("Register", mock.Anything, input).Return(resp, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockAuthService)
		handler := NewAuthHandler(mockSvc, logger)

		input := authusecase.LoginInput{
			Email:    "test@example.com",
			Password: "password123",
		}
		resp := &authusecase.AuthResponse{
			Token:    "test-token",
			UserID:   uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		mockSvc.On("Login", mock.Anything, input).Return(resp, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
