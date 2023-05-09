package servicesImplementation

import (
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/services"
)

type comparisonListServiceImplementation struct {
	comparisonListRepository repository.ComparisonListRepository
	instrumentRepository     repository.InstrumentRepository
	logger                   *logger.Logger
}

func NewComparisonListServiceImplementation(comparisonListRepository repository.ComparisonListRepository, instrumentRepository repository.InstrumentRepository, logger *logger.Logger) services.ComparisonListService {
	return &comparisonListServiceImplementation{
		comparisonListRepository: comparisonListRepository,
		instrumentRepository:     instrumentRepository,
		logger:                   logger,
	}
}

func (c *comparisonListServiceImplementation) AddInstrument(id uint64, instrumentId uint64) error {
	_, err := c.comparisonListRepository.Get(id)

	fields := map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		c.logger.WithFields(fields).Error(serviceErrors.ComparisonListAddInstrumentFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		c.logger.WithFields(fields).Error(serviceErrors.ComparisonListAddInstrumentFailed.Error() + err.Error())
		return err
	}

	_, err = c.instrumentRepository.Get(instrumentId)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		c.logger.WithFields(fields).Error(serviceErrors.ComparisonListAddInstrumentFailed.Error() + serviceErrors.InstrumentDoesNotExists.Error())
		return serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		c.logger.WithFields(fields).Error(serviceErrors.ComparisonListAddInstrumentFailed.Error() + err.Error())
		return err
	}

	err = c.comparisonListRepository.AddInstrument(id, instrumentId)
	if err != nil {
		c.logger.WithFields(fields).Error(serviceErrors.ComparisonListAddInstrumentFailed.Error() + err.Error())
		return err
	}

	c.logger.WithFields(fields).Info("comparisonList add instrument completed")
	return nil
}

func (c *comparisonListServiceImplementation) DeleteInstrument(id uint64, instrumentId uint64) error {
	_, err := c.comparisonListRepository.Get(id)

	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Error(serviceErrors.ComparisonListDeleteInstrumentFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Error(serviceErrors.ComparisonListDeleteInstrumentFailed.Error() + err.Error())
		return err
	}

	_, err = c.instrumentRepository.Get(instrumentId)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Error(serviceErrors.ComparisonListDeleteInstrumentFailed.Error() + serviceErrors.InstrumentDoesNotExists.Error())
		return serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Error(serviceErrors.ComparisonListDeleteInstrumentFailed.Error() + err.Error())
		return err
	}

	err = c.comparisonListRepository.DeleteInstrument(id, instrumentId)
	if err != nil {
		c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Error(serviceErrors.ComparisonListDeleteInstrumentFailed.Error() + err.Error())
		return err
	}
	c.logger.WithFields(map[string]interface{}{"comparisonList_id": id, "instrument_id": instrumentId}).Info("comparisonList delete instrument completed")

	return nil
}
