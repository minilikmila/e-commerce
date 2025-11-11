package auth

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/minilik/ecommerce/config"
	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
	hashpkg "github.com/minilik/ecommerce/pkg/hash"
	jwtpkg "github.com/minilik/ecommerce/pkg/jwt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	// passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^a-zA-Z0-9]).{8,}$`)
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*RegisterResponse, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
	PromoteToAdmin(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	users   repository.UserRepository
	hasher  hashpkg.Hasher
	tokens  jwtpkg.Manager
	cfg     *config.Config
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewService(
	users repository.UserRepository,
	hasher hashpkg.Hasher,
	tokens jwtpkg.Manager,
	cfg *config.Config,
	logger *zap.Logger,
) Service {
	return &service{
		users:   users,
		hasher:  hasher,
		tokens:  tokens,
		cfg:     cfg,
		logger:  logger,
		nowFunc: time.Now,
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*RegisterResponse, error) {
	if err := s.validateRegisterInput(ctx, input); err != nil {
		return nil, err
	}

	hashed, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:        uuid.New(),
		Username:  strings.TrimSpace(input.Username),
		Email:     strings.ToLower(strings.TrimSpace(input.Email)),
		Password:  hashed,
		Role:      domain.RoleUser,
		CreatedAt: s.nowFunc(),
		UpdatedAt: s.nowFunc(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(input.Password, user.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return s.issueToken(user)
}

func (s *service) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}
	if user.Role == domain.RoleAdmin {
		return nil
	}
	return s.users.UpdateRole(ctx, userID, domain.RoleAdmin)
}

func (s *service) issueToken(user *domain.User) (*AuthResponse, error) {
	ttl := s.cfg.JWT.AccessTokenTTL
	token, err := s.tokens.GenerateAccessToken(user.ID, user.Username, string(user.Role), ttl, s.cfg.JWT.Issuer)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	expiresAt := s.nowFunc().Add(ttl)

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
	}, nil
}

func (s *service) validateRegisterInput(ctx context.Context, input RegisterInput) error {
	if strings.TrimSpace(input.Username) == "" || !usernameRegex.MatchString(input.Username) {
		return domain.ErrInvalidUsernameFormat
	}

	if err := validateEmail(input.Email); err != nil {
		return err
	}

	if !isValidPassword(input.Password) {
		return domain.ErrInvalidPasswordFormat
	}

	if existing, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email))); err == nil && existing != nil {
		return domain.ErrEmailAlreadyExists
	} else if err != nil {
		return err
	}

	if existing, err := s.users.FindByUsername(ctx, strings.TrimSpace(input.Username)); err == nil && existing != nil {
		return domain.ErrUsernameAlreadyExists
	} else if err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return domain.ErrEmailCannotEmpty
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.ErrInvalidEmailFormat
	}
	return nil
}

// all registrations become regular users; admin seeding controls admin creation.

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}
