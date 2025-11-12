package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/adapter/middleware"
	"github.com/minilik/ecommerce/internal/domain"
	orderusecase "github.com/minilik/ecommerce/internal/usecase/order"
	"github.com/minilik/ecommerce/pkg/response"
)

type OrderHandler struct {
	service orderusecase.Service
	logger  *zap.Logger
}

func NewOrderHandler(service orderusecase.Service, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	// @Summary Create order
	// @Description Place a new order (user or admin)
	// @Tags Orders
	// @Accept json
	// @Produce json
	// @Param payload body orderusecase.CreateOrderInput true "Order payload"
	// @Success 201 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Security BearerAuth
	// @Router /orders [post]
	var input orderusecase.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid input", []string{err.Error()}))
		return
	}

	// read from saved context in middleware
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorBase("unauthorized", []string{"authentication required"}))
		return
	}

	if claims.Role != domain.RoleUser && claims.Role != domain.RoleAdmin {
		c.JSON(http.StatusForbidden, response.ErrorBase("forbidden", []string{"user role required"}))
		return
	}

	order, err := h.service.Create(c.Request.Context(), claims.UserID, input)
	if err != nil {
		h.logger.Warn("failed to create order", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			c.JSON(http.StatusNotFound, response.ErrorBase("product not found", []string{err.Error()}))
		case errors.Is(err, domain.ErrInsufficientStock):
			c.JSON(http.StatusBadRequest, response.ErrorBase("insufficient stock", []string{err.Error()}))
		default:
			c.JSON(http.StatusBadRequest, response.ErrorBase("failed to create order", []string{err.Error()}))
		}
		return
	}

	c.JSON(http.StatusCreated, response.SuccessBase("order created", order))
}

func (h *OrderHandler) List(c *gin.Context) {
	// @Summary List my orders
	// @Description Get current user's orders
	// @Tags Orders
	// @Produce json
	// @Success 200 {object} response.Base
	// @Security BearerAuth
	// @Router /orders [get]
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorBase("unauthorized", []string{"authentication required"}))
		return
	}

	orders, err := h.service.ListForUser(c.Request.Context(), claims.UserID)
	if err != nil {
		h.logger.Error("failed to list orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ErrorBase("failed to list orders", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessBase("orders retrieved", orders))
}
