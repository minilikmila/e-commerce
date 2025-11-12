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

	authusecase "github.com/minilik/ecommerce/internal/usecase/auth"
)

type mockAuthServiceForAdmin struct {
	mock.Mock
}

func (m *mockAuthServiceForAdmin) Register(ctx context.Context, input authusecase.RegisterInput) (*authusecase.RegisterResponse, error) {
	return nil, nil
}

func (m *mockAuthServiceForAdmin) Login(ctx context.Context, input authusecase.LoginInput) (*authusecase.AuthResponse, error) {
	return nil, nil
}

func (m *mockAuthServiceForAdmin) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestAdminHandler_PromoteUserToAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockAuthServiceForAdmin)
		handler := NewAdminHandler(mockSvc, logger)

		userID := uuid.New()
		mockSvc.On("PromoteToAdmin", mock.Anything, userID).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/"+userID.String()+"/admin", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}

		handler.PromoteUserToAdmin(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

