package repository

import "context"

// UnitOfWork defines a transaction boundary for repositories.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(tx RepositoryProvider) error) error
}

// RepositoryProvider exposes repositories bound to the same transaction context.
type RepositoryProvider interface {
	Users() UserRepository
	Products() ProductRepository
	Orders() OrderRepository
}
