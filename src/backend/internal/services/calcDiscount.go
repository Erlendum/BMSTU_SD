package services

import (
	"backend/internal/models"
)

type CalcDiscountService interface {
	CalcDiscount(user *models.User, instruments []models.Instrument) ([]models.Instrument, error)
}
