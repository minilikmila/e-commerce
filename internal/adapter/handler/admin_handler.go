package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/domain"
	authusecase "github.com/minilik/ecommerce/internal/usecase/auth"
	"github.com/minilik/ecommerce/pkg/response"
)

type AdminHandler struct {
	auth   authusecase.Service
	logger *zap.Logger
}

func NewAdminHandler(auth authusecase.Service, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{auth: auth, logger: logger}
}

// PromoteUserToAdmin promotes a user to admin (admin-only).
func (h *AdminHandler) PromoteUserToAdmin(c *gin.Context) {
	h.logger.Info("Admin promotion", zap.String("admin", ""))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid user id", []string{err.Error()}))
		return
	}
	h.logger.Info("Admin promotion", zap.String("admin", id.String()))
	if err := h.auth.PromoteToAdmin(c.Request.Context(), id); err != nil {
		if err == domain.ErrUserNotFound {
			// h.logger.Info("Admin promotion", zap.String("admin", id.String()))
			c.JSON(http.StatusNotFound, response.ErrorBase("user not found", []string{err.Error()}))
			return
		}
		h.logger.Warn("promote user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ErrorBase("failed to promote user", []string{err.Error()}))
		return
	}
	c.JSON(http.StatusOK, response.SuccessBase("user promoted to admin", nil))
}
