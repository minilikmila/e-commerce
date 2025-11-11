package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = bcrypt.DefaultCost

// Hasher defines password hashing behaviour.
type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hashed string) error
}

type bcryptHasher struct {
	cost int
}

// NewBcryptHasher returns a password hasher using the bcrypt algorithm.
func NewBcryptHasher(cost int) Hasher {
	if cost == 0 {
		cost = defaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func (b *bcryptHasher) Compare(password, hashed string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return fmt.Errorf("compare password: %w", err)
	}
	return nil
}
