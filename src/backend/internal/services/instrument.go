package services

import "backend/internal/models"

type InstrumentService interface {
	Create(instrument *models.Instrument, login string) error
	Update(id uint64, login string, fieldsToUpdate models.InstrumentFieldsToUpdate) error
	Delete(id uint64, login string) error
	Get(id uint64) (*models.Instrument, error)
	GetList() ([]models.Instrument, error)
}
