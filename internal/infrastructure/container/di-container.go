package di_container

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/minilik/ecommerce/config"
	"github.com/minilik/ecommerce/internal/adapter/handler"
	mw "github.com/minilik/ecommerce/internal/adapter/middleware"
	gormrepo "github.com/minilik/ecommerce/internal/adapter/repository/gorm"
	"github.com/minilik/ecommerce/internal/adapter/router"
	"github.com/minilik/ecommerce/internal/infrastructure/database"
	authusecase "github.com/minilik/ecommerce/internal/usecase/auth"
	orderusecase "github.com/minilik/ecommerce/internal/usecase/order"
	productusecase "github.com/minilik/ecommerce/internal/usecase/product"
	"github.com/minilik/ecommerce/pkg/cache"
	"github.com/minilik/ecommerce/pkg/cloudinary"
	hashpkg "github.com/minilik/ecommerce/pkg/hash"
	jwtpkg "github.com/minilik/ecommerce/pkg/jwt"
	"github.com/minilik/ecommerce/pkg/logger"
)

type DIContainer struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Router *gin.Engine
}

// Build initializes and wires all application dependencies... DI container pattern
func Build(cfg *config.Config) (*DIContainer, error) {
	log, err := logger.New(cfg.App.Environment)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	db, err := database.NewPostgres(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("create database connection: %w", err)
	}

	if err := database.Migrate(db); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	hasher := hashpkg.NewBcryptHasher(0)
	jwtManager, err := jwtpkg.NewManager(cfg.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("create jwt manager: %w", err)
	}

	userRepo := gormrepo.NewUserRepository(db)
	productRepo := gormrepo.NewProductRepository(db)
	orderRepo := gormrepo.NewOrderRepository(db)
	uow := gormrepo.NewUnitOfWork(db)

	authService := authusecase.NewService(userRepo, hasher, jwtManager, cfg, log)
	var prodCache *cache.MemoryCache
	if cfg.Cache.Enabled {
		prodCache = cache.NewMemoryCache(cfg.Cache.ProductListTTL, cfg.Cache.MaxProductEntries)
	}
	productService := productusecase.NewService(productRepo, orderRepo, log, prodCache)
	orderService := orderusecase.NewService(uow, log)

	// Cloudinary uploader + image repo/service
	var uploader *cloudinary.Client
	if cfg.Cloud.CloudName != "" && (cfg.Cloud.UploadPreset != "" || cfg.Cloud.APIKey != "") {
		uploader = cloudinary.NewClient(cfg.Cloud.CloudName, cfg.Cloud.APIKey, cfg.Cloud.APISecret, cfg.Cloud.UploadPreset, cfg.Cloud.Folder)
	}
	imageRepo := gormrepo.NewProductImageRepository(db)
	imageService := productusecase.NewImageService(imageRepo, uploader, log)

	authHandler := handler.NewAuthHandler(authService, log)
	productHandler := handler.NewProductHandler(productService, log).WithImageService(imageService)
	orderHandler := handler.NewOrderHandler(orderService, log)

	authMiddleware := mw.NewAuthMiddleware(log, jwtManager)
	var rateLimiter *mw.RateLimitMiddleware
	if cfg.Rate.Enabled && cfg.Rate.Limit > 0 && cfg.Rate.Window > 0 {
		rateLimiter = mw.NewRateLimitMiddleware(cfg.Rate.Limit, cfg.Rate.Window)
	}

	engine := router.Setup(router.Dependencies{
		AuthHandler:    authHandler,
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
		AuthMiddleware: authMiddleware,
		RateLimiter:    rateLimiter,
	})

	return &DIContainer{
		Config: cfg,
		Logger: log,
		DB:     db,
		Router: engine,
	}, nil
}

// Close releases resources held by the container.
func (c *DIContainer) Close() error {
	logger.Sync(c.Logger)
	if c.DB == nil {
		return nil
	}
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
