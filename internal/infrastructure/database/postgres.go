package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/minilik/ecommerce/config"
	"github.com/minilik/ecommerce/internal/adapter/repository/gorm/models"
)

// NewPostgres creates a new gorm.DB instance connected to PostgreSQL.
func NewPostgres(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	gormLogger := logger.New(
		zap.NewStdLog(log),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get generic db: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Migrate runs database migrations.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
}
