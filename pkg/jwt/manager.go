package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	userIDClaimKey   = "uid"
	usernameClaimKey = "uname"
	roleClaimKey     = "role"
)

// Claims represents the JWT claims used by the application.
type Claims struct {
	UserID   uuid.UUID
	Username string
	Role     string
	jwt.RegisteredClaims
}

// Manager defines operations for generating and validating JWT tokens.
type Manager interface {
	GenerateAccessToken(userID uuid.UUID, username, role string, ttl time.Duration, issuer string) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

type manager struct {
	secret []byte
}

// NewManager creates a new JWT manager with the provided secret.
func NewManager(secret string) (Manager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret cannot be empty")
	}

	return &manager{
		secret: []byte(secret),
	}, nil
}

func (m *manager) GenerateAccessToken(userID uuid.UUID, username, role string, ttl time.Duration, issuer string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		userIDClaimKey:   userID.String(),
		usernameClaimKey: username,
		roleClaimKey:     role,
		"iss":            issuer,
		"iat":            now.Unix(),
		"exp":            now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return str, nil
}

func (m *manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := mapClaims[userIDClaimKey].(string)
	if !ok {
		return nil, errors.New("user id claim missing")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id claim: %w", err)
	}

	username, _ := mapClaims[usernameClaimKey].(string)
	role, _ := mapClaims[roleClaimKey].(string)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
	}

	if iss, ok := mapClaims["iss"].(string); ok {
		claims.Issuer = iss
	}
	if exp, ok := mapClaims["exp"].(float64); ok {
		claims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(exp), 0))
	}
	if iat, ok := mapClaims["iat"].(float64); ok {
		claims.IssuedAt = jwt.NewNumericDate(time.Unix(int64(iat), 0))
	}

	return claims, nil
}
