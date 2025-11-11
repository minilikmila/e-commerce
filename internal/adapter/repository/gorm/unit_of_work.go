package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/minilik/ecommerce/internal/domain/repository"
)

type unitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork returns a UnitOfWork implementation backed by GORM transactions.
func NewUnitOfWork(db *gorm.DB) repository.UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Execute(ctx context.Context, fn func(tx repository.RepositoryProvider) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		provider := &repositoryProvider{
			users:    NewUserRepository(tx),
			products: NewProductRepository(tx),
			orders:   NewOrderRepository(tx),
		}
		return fn(provider)
	})
}

type repositoryProvider struct {
	users    repository.UserRepository
	products repository.ProductRepository
	orders   repository.OrderRepository
}

func (p *repositoryProvider) Users() repository.UserRepository {
	return p.users
}

func (p *repositoryProvider) Products() repository.ProductRepository {
	return p.products
}

func (p *repositoryProvider) Orders() repository.OrderRepository {
	return p.orders
}
