package repository

import "backend/internal/models"

type ComparisonListRepository interface {
	Create(comparisonList *models.ComparisonList) error
	AddInstrument(id uint64, instrumentId uint64) error
	Update(id uint64, fieldsToUpdate models.ComparisonListFieldsToUpdate) error
	DeleteInstrument(id uint64, instrumentId uint64) error
	Get(userId uint64) (*models.ComparisonList, error)
	GetUser(id uint64) (*models.User, error)
	GetInstruments(userId uint64) ([]models.Instrument, error)
}
