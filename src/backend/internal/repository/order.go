package repository

import "backend/internal/models"

type OrderRepository interface {
	Create(order *models.Order) (uint64, error)
	CreateOrderElement(element *models.OrderElement) error
}
