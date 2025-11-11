package response

import (
	"math"
	"net/http"
)

// Base represents the standard API response body.
type Base struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// Paginated represents a paginated API response body.
type Paginated struct {
	Success       bool        `json:"success"`
	Message       string      `json:"message"`
	Data          interface{} `json:"data"`
	CurrentPage   int         `json:"currentPage"`
	PageSize      int         `json:"pageSize"`
	TotalPages    int         `json:"totalPages"`
	TotalProducts int64       `json:"totalProducts"`
	Errors        []string    `json:"errors,omitempty"`
}

// SuccessBase returns a successful base response.
func SuccessBase(message string, object interface{}) Base {
	return Base{
		Success: true,
		Message: message,
		Data:    object,
	}
}

// ErrorBase returns an error base response.
// TODO: check the env't and ignore error details if it's in production
func ErrorBase(message string, errs []string) Base {
	return Base{
		Success: false,
		Message: message,
		Errors:  errs,
	}
}

// SuccessPaginated returns a successful paginated response.
func SuccessPaginated(message string, object interface{}, page, size int, total int64) Paginated {
	totalPages := 0
	if size > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(size)))
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return Paginated{
		Success:       true,
		Message:       message,
		Data:          object,
		CurrentPage:   page,
		PageSize:      size,
		TotalPages:    totalPages,
		TotalProducts: total,
	}
}

// StatusCodeFromBool returns http status based on success.
func StatusCodeFromBool(success bool) int {
	if success {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
