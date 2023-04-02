package servicesImplementation

import (
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/repository"
	"backend/internal/services"
)

type comparisonListServiceImplementation struct {
	comparisonListRepository repository.ComparisonListRepository
	instrumentRepository     repository.InstrumentRepository
}

func NewComparisonListServiceImplementation(comparisonListRepository repository.ComparisonListRepository, instrumentRepository repository.InstrumentRepository) services.ComparisonListService {
	return &comparisonListServiceImplementation{
		comparisonListRepository: comparisonListRepository,
		instrumentRepository:     instrumentRepository,
	}
}

func (c *comparisonListServiceImplementation) AddInstrument(id uint64, instrumentId uint64) error {
	_, err := c.comparisonListRepository.Get(id)

	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		return err
	}

	_, err = c.instrumentRepository.Get(instrumentId)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		return err
	}

	return c.comparisonListRepository.AddInstrument(id, instrumentId)
}

func (c *comparisonListServiceImplementation) DeleteInstrument(id uint64, instrumentId uint64) error {
	_, err := c.comparisonListRepository.Get(id)

	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		return err
	}

	_, err = c.instrumentRepository.Get(instrumentId)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		return err
	}

	return c.comparisonListRepository.DeleteInstrument(id, instrumentId)
}
