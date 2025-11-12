package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/minilik/ecommerce/internal/adapter/handler"
	"github.com/minilik/ecommerce/internal/adapter/middleware"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/pkg/response"

	// Import usecase packages for Swagger type references
	_ "github.com/minilik/ecommerce/internal/usecase/auth"
	_ "github.com/minilik/ecommerce/internal/usecase/order"
	_ "github.com/minilik/ecommerce/internal/usecase/product"

	// Import docs for Swagger
	_ "github.com/minilik/ecommerce/docs"
)

const (
	APIBasePath = "/api/v1"
)

type Dependencies struct {
	AuthHandler    *handler.AuthHandler
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
	AdminHandler   *handler.AdminHandler
	AuthMiddleware *middleware.AuthMiddleware
	RateLimiter    *middleware.RateLimitMiddleware
}

// COMMENTS ARE FOR SWAGGER DOCS PURPOSES TO ENABLE AUTOMATICALLY GENERATING THE DOCS FROM THE CODE

// @title E-commerce API
// @version 1.0
// @description REST API for an E-commerce platform with authentication, products, orders, and image uploads.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func Setup(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CorsMiddleware())

	// Swagger UI - register before rate limiter to exclude it
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply rate limiter only to API routes (excludes Swagger)
	if deps.RateLimiter != nil {
		r.Use(func(c *gin.Context) {
			// Skip rate limiting for Swagger routes
			if c.Request.URL.Path == "/swagger" || len(c.Request.URL.Path) > 8 && c.Request.URL.Path[:9] == "/swagger/" {
				c.Next()
				return
			}
			deps.RateLimiter.RateLimit()(c)
		})
	}

	v1 := r.Group(APIBasePath) // versioning apis
	v1.GET("/health", func(c *gin.Context) {
		// @Summary Health check
		// @Description Check API health status
		// @Tags Health
		// @Produce json
		// @Success 200 {object} response.Base
		// @Router /health [get]
		c.JSON(200, response.SuccessBase("ok", nil))
	})
	// auth endpoints: public access
	auth := v1.Group("/auth")
	{
		// @Summary Register a new user
		// @Description Create a new user account (role=user)
		// @Tags Auth
		// @Accept json
		// @Produce json
		// @Param payload body authusecase.RegisterInput true "Register payload"
		// @Success 201 {object} response.Base
		// @Failure 400 {object} response.Base
		// @Router /auth/register [post]
		auth.POST("/register", deps.AuthHandler.Register)

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
		auth.POST("/login", deps.AuthHandler.Login)
	}
	// Query endpoints: Public access
	product := v1.Group("/products")
	{
		// @Summary List products
		// @Description List products with pagination (public)
		// @Tags Products
		// @Produce json
		// @Param page query int false "Page number"
		// @Param limit query int false "Page size"
		// @Param search query string false "Search term"
		// @Success 200 {object} response.Paginated
		// @Router /products [get]
		product.GET("", deps.ProductHandler.List)

		// @Summary Get product
		// @Description Get product details (public)
		// @Tags Products
		// @Produce json
		// @Param id path string true "Product ID"
		// @Success 200 {object} response.Base
		// @Failure 404 {object} response.Base
		// @Router /products/{id} [get]
		product.GET("/:id", deps.ProductHandler.Get)
	}
	// Mutation endpoints for admin
	adminProducts := v1.Group("/products")
	adminProducts.Use(deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireRoles(domain.RoleAdmin))
	{
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
		adminProducts.POST("", deps.ProductHandler.Create)

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
		adminProducts.PUT("/:id", deps.ProductHandler.Update)

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
		adminProducts.DELETE("/:id", deps.ProductHandler.Delete)

		// @Summary Upload product images
		// @Description Upload up to 4 images for a product (admin only)
		// @Tags Products
		// @Accept multipart/form-data
		// @Produce json
		// @Param id path string true "Product ID"
		// @Param files formData file true "Image files"
		// @Success 201 {object} response.Base
		// @Security BearerAuth
		// @Router /products/{id}/images [post]
		adminProducts.POST("/:id/images", deps.ProductHandler.UploadImages)
	}

	// Mutation endpoints for user and admin role
	orders := v1.Group("/orders")
	orders.Use(deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireRoles(domain.RoleAdmin, domain.RoleUser))
	{
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
		orders.POST("", deps.OrderHandler.Create)

		// @Summary List my orders
		// @Description Get current user's orders
		// @Tags Orders
		// @Produce json
		// @Success 200 {object} response.Base
		// @Security BearerAuth
		// @Router /orders [get]
		orders.GET("", deps.OrderHandler.List)
	}

	// Admin endpoints
	admin := v1.Group("/admin")
	admin.Use(deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireRoles(domain.RoleAdmin))
	{
		// @Summary Promote user to admin
		// @Description Promote a user to admin role (admin only)
		// @Tags Admin
		// @Produce json
		// @Param id path string true "User ID"
		// @Success 200 {object} response.Base
		// @Failure 404 {object} response.Base
		// @Security BearerAuth
		// @Router /admin/users/{id}/admin [post]
		admin.POST("/users/:id/admin", deps.AdminHandler.PromoteUserToAdmin)
	}

	return r
}
