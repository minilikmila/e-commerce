package domain

import "errors"

var (
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrUsernameAlreadyExists   = errors.New("username already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrProductNotFound         = errors.New("product not found")
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrInvalidPasswordFormat   = errors.New("invalid password format")
	ErrInvalidUsernameFormat   = errors.New("invalid username format: username must be alphanumeric without spaces")
	ErrInvalidEmailFormat      = errors.New("invalid email format")
	ErrEmailCannotEmpty        = errors.New("email cannot be empty")
	ErrProductHasPendingOrders = errors.New("cannot delete product: product has pending orders")
)
