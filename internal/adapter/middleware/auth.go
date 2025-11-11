package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/domain"
	jwtpkg "github.com/minilik/ecommerce/pkg/jwt"
	"github.com/minilik/ecommerce/pkg/response"
)

const userContextKey = "currentUser"

type UserClaims struct {
	UserID   uuid.UUID
	Username string
	Role     domain.Role
}

type AuthMiddleware struct {
	logger *zap.Logger
	jwt    jwtpkg.Manager
}

func NewAuthMiddleware(logger *zap.Logger, jwt jwtpkg.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
		jwt:    jwt,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorBase("authorization token missing", []string{"authorization header missing"}))
			c.Abort()
			return
		}
		// TODO: add read token from cookie if not in header

		claims, err := a.jwt.ParseToken(token)
		if err != nil {
			a.logger.Warn("failed to parse token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, response.ErrorBase("invalid token", []string{err.Error()}))
			c.Abort()
			return
		}

		userClaims := UserClaims{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     domain.Role(claims.Role),
		}

		c.Set(userContextKey, userClaims)
		c.Next()
	}
}

func (a *AuthMiddleware) RequireRoles(roles ...domain.Role) gin.HandlerFunc {
	roleSet := make(map[domain.Role]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(userContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, response.ErrorBase("unauthorized", []string{"authentication required"}))
			c.Abort()
			return
		}

		claims, ok := value.(UserClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorBase("unauthorized", []string{"invalid authentication context"}))
			c.Abort()
			return
		}

		if len(roleSet) == 0 {
			c.Next()
			return
		}

		if _, ok := roleSet[claims.Role]; !ok {
			c.JSON(http.StatusForbidden, response.ErrorBase("forbidden", []string{"insufficient permissions"}))
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserClaims(c *gin.Context) (UserClaims, bool) {
	value, exists := c.Get(userContextKey)
	if !exists {
		return UserClaims{}, false
	}

	claims, ok := value.(UserClaims)
	return claims, ok
}

func extractToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
