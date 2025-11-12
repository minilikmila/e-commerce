package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/adapter/middleware"
	"github.com/minilik/ecommerce/internal/domain"
	productusecase "github.com/minilik/ecommerce/internal/usecase/product"
	"github.com/minilik/ecommerce/pkg/response"
)

type ProductHandler struct {
	service      productusecase.Service
	imageService productusecase.ImageService
	logger       *zap.Logger
}

func NewProductHandler(service productusecase.Service, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) WithImageService(img productusecase.ImageService) *ProductHandler {
	h.imageService = img
	return h
}

func (h *ProductHandler) Create(c *gin.Context) {
	// @Summary Create product
	// @Description Create a product (admin only)
	// @Tags Products
	// @Accept json
	// @Produce json
	// @Param payload body productusecase.CreateProductInput true "Product payload"
	// @Success 201 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Security BearerAuth
	// @Router /products [post]
	var input productusecase.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid input", []string{err.Error()}))
		return
	}
	// read from saved context in middleware
	claims, ok := middleware.GetUserClaims(c)
	if !ok || claims.Role != domain.RoleAdmin {
		c.JSON(http.StatusForbidden, response.ErrorBase("forbidden", []string{"admin role required"}))
		return
	}

	product, err := h.service.Create(c.Request.Context(), claims.UserID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("failed to create product", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessBase("product created", product))
}

func (h *ProductHandler) Update(c *gin.Context) {
	// @Summary Update product
	// @Description Update product fields (admin only)
	// @Tags Products
	// @Accept json
	// @Produce json
	// @Param id path string true "Product ID"
	// @Param payload body productusecase.UpdateProductInput true "Update payload"
	// @Success 200 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Failure 404 {object} response.Base
	// @Security BearerAuth
	// @Router /products/{id} [put]
	var input productusecase.UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid input", []string{err.Error()}))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid product id", []string{err.Error()}))
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.ErrorBase("product not found", []string{err.Error()}))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorBase("failed to update product", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessBase("product updated", product))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	// @Summary Delete product
	// @Description Delete a product if no pending orders (admin only)
	// @Tags Products
	// @Produce json
	// @Param id path string true "Product ID"
	// @Success 200 {object} response.Base
	// @Failure 400 {object} response.Base
	// @Failure 404 {object} response.Base
	// @Security BearerAuth
	// @Router /products/{id} [delete]
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid product id", []string{err.Error()}))
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.ErrorBase("product not found", []string{err.Error()}))
			return
		}
		if err == domain.ErrProductHasPendingOrders {
			c.JSON(http.StatusBadRequest, response.ErrorBase("cannot delete product", []string{err.Error()}))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorBase("failed to delete product", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessBase("product deleted", nil))
}

func (h *ProductHandler) Get(c *gin.Context) {
	// @Summary Get product
	// @Description Get product details (public)
	// @Tags Products
	// @Produce json
	// @Param id path string true "Product ID"
	// @Success 200 {object} response.Base
	// @Failure 404 {object} response.Base
	// @Router /products/{id} [get]
	// this is also allowed for public access
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid product id", []string{err.Error()}))
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.ErrorBase("product not found", []string{err.Error()}))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorBase("failed to fetch product", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessBase("product retrieved", product))
}

func (h *ProductHandler) List(c *gin.Context) {
	// @Summary List products
	// @Description List products with pagination (public)
	// @Tags Products
	// @Produce json
	// @Param page query int false "Page number"
	// @Param limit query int false "Page size"
	// @Param search query string false "Search term"
	// @Success 200 {object} response.Paginated
	// @Router /products [get]
	// this is also allowed for public access : it returns list of products
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "limit", 10)
	search := c.Query("search")

	products, total, err := h.service.List(c.Request.Context(), productusecase.ListProductsInput{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		h.logger.Error("failed to list products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ErrorBase("failed to list products", []string{err.Error()}))
		return
	}

	resp := response.SuccessPaginated(
		"products retrieved",
		products,
		page,
		pageSize,
		total,
	)

	c.JSON(http.StatusOK, resp)
}

func parseQueryInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func (h *ProductHandler) UploadImages(c *gin.Context) {
	// @Summary Upload product images
	// @Description Upload up to 4 images for a product (admin only)
	// @Tags Products
	// @Accept multipart/form-data
	// @Produce json
	// @Param id path string true "Product ID"
	// @Param files formData file true "Image files" collectionFormat(multi)
	// @Success 201 {object} response.Base
	// @Security BearerAuth
	// @Router /products/{id}/images [post]
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid product id", []string{err.Error()}))
		return
	}
	if h.imageService == nil {
		c.JSON(http.StatusInternalServerError, response.ErrorBase("image service not configured", []string{}))
		return
	}
	claims, ok := middleware.GetUserClaims(c)
	if !ok || claims.Role != domain.RoleAdmin {
		c.JSON(http.StatusForbidden, response.ErrorBase("forbidden", []string{"admin role required"}))
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("invalid multipart form", []string{err.Error()}))
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, response.ErrorBase("no files uploaded", []string{}))
		return
	}
	if len(files) > 4 {
		c.JSON(http.StatusBadRequest, response.ErrorBase("maximum 4 images allowed", []string{}))
		return
	}
	uploaded, err := h.imageService.UploadImages(c.Request.Context(), id, files)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBase("failed to upload images", []string{err.Error()}))
		return
	}
	c.JSON(http.StatusCreated, response.SuccessBase("images uploaded", uploaded))
}
