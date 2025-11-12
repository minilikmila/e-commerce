package router

// Import packages for Swagger type references in annotations
// These are only used in Swagger annotations, not in code

// This file helps Swaggo find route annotations
// Swaggo has difficulty scanning methods on structs, so we document routes here
// These dummy functions with annotations ensure Swaggo can generate the docs

// @Summary Register a new user
// @Description Create a new user account (role=user)
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body auth.RegisterInput true "Register payload"
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Router /auth/register [post]
func _() {}

// @Summary Login
// @Description Authenticate and obtain JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body auth.LoginInput true "Login payload"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/login [post]
func _() {}

// @Summary List products
// @Description List products with pagination (public)
// @Tags Products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search term"
// @Success 200 {object} response.Paginated
// @Router /products [get]
func _() {}

// @Summary Get product
// @Description Get product details (public)
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Base
// @Failure 404 {object} response.Base
// @Router /products/{id} [get]
func _() {}

// @Summary Create product
// @Description Create a product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param payload body product.CreateProductInput true "Product payload"
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Security BearerAuth
// @Router /products [post]
func _() {}

// @Summary Update product
// @Description Update product fields (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param payload body product.UpdateProductInput true "Update payload"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Security BearerAuth
// @Router /products/{id} [put]
func _() {}

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
func _() {}

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
func _() {}

// @Summary Create order
// @Description Place a new order (user or admin)
// @Tags Orders
// @Accept json
// @Produce json
// @Param payload body order.CreateOrderInput true "Order payload"
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Security BearerAuth
// @Router /orders [post]
func _() {}

// @Summary List my orders
// @Description Get current user's orders
// @Tags Orders
// @Produce json
// @Success 200 {object} response.Base
// @Security BearerAuth
// @Router /orders [get]
func _() {}

// @Summary Promote user to admin
// @Description Promote a user to admin role (admin only)
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Base
// @Failure 404 {object} response.Base
// @Security BearerAuth
// @Router /admin/users/{id}/admin [post]
func _() {}
