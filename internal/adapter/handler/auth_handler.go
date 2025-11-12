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
	// @Summary Register a new user
	// @Description Create a new user account (role=user)
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param payload body authusecase.RegisterInput true "Register payload"
	// @Success 201 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Router /auth/register [post]
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
	// @Summary Login
	// @Description Authenticate and obtain JWT token
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param payload body authusecase.LoginInput true "Login payload"
	// @Success 200 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Failure 401 {object} response.Base
	// @Router /auth/login [post]
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
