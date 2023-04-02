package services

import "backend/internal/models"

type UserService interface {
	Create(user *models.User, password string) error
	Get(login, password string) (*models.User, error)
	GetComparisonList(id uint64) (*models.ComparisonList, []models.Instrument, error)
}
