package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/repository"
	"backend/internal/services"
)

type instrumentServiceImplementation struct {
	instrumentRepository repository.InstrumentRepository
	userRepository       repository.UserRepository
}

func NewInstrumentServiceImplementation(instrumentRepository repository.InstrumentRepository, userRepository repository.UserRepository) services.InstrumentService {
	return &instrumentServiceImplementation{
		instrumentRepository: instrumentRepository,
		userRepository:       userRepository,
	}
}

func (i *instrumentServiceImplementation) Create(instrument *models.Instrument, login string) error {
	user, err := i.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.UserDoesNotExists
	} else if err != nil {
		return err
	}

	if !user.IsAdmin {
		return serviceErrors.UserCanNotCreateInstrument
	}

	return i.instrumentRepository.Create(instrument)
}

func (i *instrumentServiceImplementation) Update(id uint64, login string, fieldsToUpdate models.InstrumentFieldsToUpdate) error {
	canUpdate, err := i.checkCanUserChangeInstrument(id, login)
	if err != nil {
		return err
	} else if !canUpdate {
		return serviceErrors.UserCanNotUpdateInstrument
	}

	return i.instrumentRepository.Update(id, fieldsToUpdate)
}

func (i *instrumentServiceImplementation) Delete(id uint64, login string) error {
	canDelete, err := i.checkCanUserChangeInstrument(id, login)
	if err != nil {
		return err
	} else if !canDelete {
		return serviceErrors.UserCanNotDeleteInstrument
	}

	return i.instrumentRepository.Delete(id)
}

func (i *instrumentServiceImplementation) Get(id uint64) (*models.Instrument, error) {
	instrument, err := i.instrumentRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, serviceErrors.InstrumentDoesNotExists
	} else if err != nil {
		return nil, err
	}

	return instrument, nil
}

func (i *instrumentServiceImplementation) GetList() ([]models.Instrument, error) {
	instruments, err := i.instrumentRepository.GetList()
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, serviceErrors.InstrumentsDoesNotExists
	} else if err != nil {
		return nil, err
	}

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
