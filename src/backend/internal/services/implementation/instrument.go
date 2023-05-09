package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/services"
)

type instrumentServiceImplementation struct {
	instrumentRepository repository.InstrumentRepository
	userRepository       repository.UserRepository
	logger               *logger.Logger
}

func NewInstrumentServiceImplementation(instrumentRepository repository.InstrumentRepository, userRepository repository.UserRepository, logger *logger.Logger) services.InstrumentService {
	return &instrumentServiceImplementation{
		instrumentRepository: instrumentRepository,
		userRepository:       userRepository,
		logger:               logger,
	}
}

func (i *instrumentServiceImplementation) Create(instrument *models.Instrument, login string) error {
	fields := map[string]interface{}{"user_login": login, "instrument_name": instrument.Name}
	user, err := i.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentCreateFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return serviceErrors.UserDoesNotExists
	} else if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentCreateFailed.Error() + err.Error())
		return err
	}

	if !user.IsAdmin {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentCreateFailed.Error() + serviceErrors.UserCanNotCreateInstrument.Error())
		return serviceErrors.UserCanNotCreateInstrument
	}

	err = i.instrumentRepository.Create(instrument)
	if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentCreateFailed.Error() + err.Error())
		return err
	}

	i.logger.WithFields(fields).Info("instrument create completed")

	return nil
}

func (i *instrumentServiceImplementation) Update(id uint64, login string, fieldsToUpdate models.InstrumentFieldsToUpdate) error {
	fields := map[string]interface{}{"user_login": login, "instrument_id": id}
	canUpdate, err := i.checkCanUserChangeInstrument(id, login)
	if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentUpdateFailed.Error() + err.Error())
		return err
	} else if !canUpdate {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentUpdateFailed.Error() + serviceErrors.UserCanNotUpdateInstrument.Error())
		return serviceErrors.UserCanNotUpdateInstrument
	}

	err = i.instrumentRepository.Update(id, fieldsToUpdate)
	if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentUpdateFailed.Error() + err.Error())
		return err
	}

	i.logger.WithFields(fields).Info("instrument update completed")

	return nil
}

func (i *instrumentServiceImplementation) Delete(id uint64, login string) error {
	fields := map[string]interface{}{"user_login": login, "instrument_id": id}
	canDelete, err := i.checkCanUserChangeInstrument(id, login)
	if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentDeleteFailed.Error() + err.Error())
		return err
	} else if !canDelete {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentDeleteFailed.Error() + serviceErrors.UserCanNotDeleteInstrument.Error())
		return serviceErrors.UserCanNotDeleteInstrument
	}

	err = i.instrumentRepository.Delete(id)
	if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentDeleteFailed.Error() + err.Error())
		return err
	}

	i.logger.WithFields(fields).Info("instrument delete completed")

	return nil
}

func (i *instrumentServiceImplementation) Get(id uint64) (*models.Instrument, error) {
	fields := map[string]interface{}{"instrument_id": id}
	instrument, err := i.instrumentRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentGetFailed.Error() + serviceErrors.InstrumentDoesNotExists.Error())
		return nil, serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		i.logger.WithFields(fields).Error(serviceErrors.InstrumentGetFailed.Error() + err.Error())
		return nil, err
	}

	i.logger.WithFields(fields).Info("instrument get completed")

	return instrument, nil
}

func (i *instrumentServiceImplementation) GetList() ([]models.Instrument, error) {
	instruments, err := i.instrumentRepository.GetList()
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		i.logger.Error(serviceErrors.InstrumentsListGetFailed.Error() + serviceErrors.InstrumentsDoesNotExists.Error())
		return nil, serviceErrors.InstrumentsDoesNotExists
	} else if err != nil {
		i.logger.Error(serviceErrors.InstrumentsListGetFailed.Error() + err.Error())
		return nil, err
	}

	i.logger.Info("instrument get list completed")
	return instruments, nil
}

func (i *instrumentServiceImplementation) checkCanUserChangeInstrument(instrumentId uint64, login string) (bool, error) {
	_, err := i.instrumentRepository.Get(instrumentId)
	if err == repositoryErrors.ObjectDoesNotExists {
		return false, serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		return false, err
	}

	user, err := i.userRepository.Get(login)
	if err == repositoryErrors.ObjectDoesNotExists {
		return false, serviceErrors.UserDoesNotExists
	} else if err != nil {
		return false, err
	}

	if user.IsAdmin {
		return true, nil
	}
	return false, nil
}
