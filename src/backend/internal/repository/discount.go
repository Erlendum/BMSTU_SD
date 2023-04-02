package repository

import "backend/internal/models"

type DiscountRepository interface {
	Create(discount *models.Discount) error
	Update(id uint64, fieldsToUpdate models.DiscountFieldsToUpdate) error
	Delete(id uint64) error
	Get(id uint64) (*models.Discount, error)
	GetList() ([]models.Discount, error)
	GetSpecificList(instrumentId uint64, userId uint64) ([]models.Discount, error)
}
