package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/internal/domain"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:100;not null"`
	Email     string    `gorm:"uniqueIndex;size:255;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"size:20;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product
	Orders   []Order
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      domain.Role(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func UserFromDomain(user *domain.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
