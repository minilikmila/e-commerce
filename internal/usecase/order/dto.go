package order

import "github.com/google/uuid"

type OrderItemInput struct {
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
}

type CreateOrderInput struct {
	Description string           `json:"description"`
	Items       []OrderItemInput `json:"items"`
}
