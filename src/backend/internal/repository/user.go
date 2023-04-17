package repository

import "backend/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	Get(login string) (*models.User, error)
	GetById(id uint64) (*models.User, error)
	GetList() ([]models.User, error)
}
