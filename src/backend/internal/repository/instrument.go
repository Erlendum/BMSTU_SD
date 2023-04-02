package repository

import "backend/internal/models"

type InstrumentRepository interface {
	Create(instrument *models.Instrument) error
	Update(id uint64, fieldsToUpdate models.InstrumentFieldsToUpdate) error
	Delete(id uint64) error
	Get(id uint64) (*models.Instrument, error)
	GetList() ([]models.Instrument, error)
}
