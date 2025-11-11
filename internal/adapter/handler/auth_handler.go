package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/domain"
	authusecase "github.com/minilik/ecommerce/internal/usecase/auth"
	"github.com/minilik/ecommerce/pkg/response"
)

type AuthHandler struct {
	service authusecase.Service
	logger  *zap.Logger
}

func NewAuthHandler(service authusecase.Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input authusecase.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid input", []string{err.Error()}))
		return
	}

	res, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists, domain.ErrUsernameAlreadyExists, domain.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, response.ErrorBase(err.Error(), []string{err.Error()}))
		default:
			h.logger.Error("register failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.ErrorBase("registration failed", []string{err.Error()}))
		}
		return
	}

	c.JSON(http.StatusCreated, response.SuccessBase("user registered successfully", res))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input authusecase.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid input", []string{err.Error()}))
		return
	}

	res, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, response.ErrorBase("invalid credentials", []string{err.Error()}))
			return
		}
		h.logger.Error("login failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ErrorBase("login failed", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessBase("login successful", res))
}
