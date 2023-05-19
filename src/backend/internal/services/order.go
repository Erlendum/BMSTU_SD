package services

import "backend/internal/models"

type OrderService interface {
	Create(userId uint64) (uint64, error)
	GetList(userId uint64) ([]models.Order, error)
	GetListForAll() ([]models.Order, error)
	Update(id uint64, login string, fieldsToUpdate models.OrderFieldsToUpdate) error
	GetOrderElements(id uint64) ([]models.OrderElement, error)
}
