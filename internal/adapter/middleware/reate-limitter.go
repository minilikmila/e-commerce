package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware handles rate limiting : we don't use redis or other external services for this just for simplicity we keep it in memory
//
//	Here recommend to use centralized rate limiting service to handle rate limiting for production environment / like in distributed system setup
type RateLimitMiddleware struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(limit int, window time.Duration) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit middleware that limits requests per IP
func (m *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		m.mutex.Lock()
		defer m.mutex.Unlock()

		now := time.Now()
		windowStart := now.Add(-m.window)

		// Clean old requests
		if requests, exists := m.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(windowStart) {
					validRequests = append(validRequests, reqTime)
				}
			}
			m.requests[clientIP] = validRequests
		}

		// Check if limit exceeded
		if len(m.requests[clientIP]) >= m.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		// Add current request
		m.requests[clientIP] = append(m.requests[clientIP], now)

		c.Next()
	}
}
