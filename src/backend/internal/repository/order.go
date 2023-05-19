package repository

import "backend/internal/models"

type OrderRepository interface {
	Create(order *models.Order) (uint64, error)
	CreateOrderElement(element *models.OrderElement) error
	GetList(userId uint64) ([]models.Order, error)
	GetListForAll() ([]models.Order, error)
	Update(id uint64, fieldsToUpdate models.OrderFieldsToUpdate) error
	GetOrderElements(id uint64) ([]models.OrderElement, error)
}
