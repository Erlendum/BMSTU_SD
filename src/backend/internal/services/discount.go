package services

import "backend/internal/models"

type DiscountService interface {
	Create(discount *models.Discount, login string) error
	CreateForAll(discount *models.Discount, login string) error
	Update(id uint64, login string, fieldsToUpdate models.DiscountFieldsToUpdate) error
	Delete(id uint64, login string) error
	Get(id uint64) (*models.Discount, error)
	GetList() ([]models.Discount, error)
}
