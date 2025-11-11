package router

import (
	"github.com/gin-gonic/gin"

	"github.com/minilik/ecommerce/internal/adapter/handler"
	"github.com/minilik/ecommerce/internal/adapter/middleware"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/pkg/response"
)

const (
	APIBasePath = "/api/v1"
)

type Dependencies struct {
	AuthHandler    *handler.AuthHandler
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func Setup(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CorsMiddleware())

	v1 := r.Group(APIBasePath) // versioning apis
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.SuccessBase("ok", nil))
	})
	// auth endpoints: public access
	auth := v1.Group("/auth")
	{
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/login", deps.AuthHandler.Login)
	}
	// Query endpoints: Public access
	product := v1.Group("/products")
	{
		product.GET("", deps.ProductHandler.List)
		product.GET("/:id", deps.ProductHandler.Get)
	}
	// Mutation endpoints for admin
	adminProducts := v1.Group("/products")
	adminProducts.Use(deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireRoles(domain.RoleAdmin))
	{
		adminProducts.POST("", deps.ProductHandler.Create)
		adminProducts.PUT("/:id", deps.ProductHandler.Update)
		adminProducts.DELETE("/:id", deps.ProductHandler.Delete)
	}

	// Mutation endpoints for user and admin role
	orders := v1.Group("/orders")
	orders.Use(deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireRoles(domain.RoleAdmin, domain.RoleUser))
	{
		orders.POST("", deps.OrderHandler.Create)
		orders.GET("", deps.OrderHandler.List)
	}

	return r
}
