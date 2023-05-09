package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/hasher"
	"backend/internal/pkg/hasher/implementation"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/services"
)

type userServiceImplementation struct {
	userRepository           repository.UserRepository
	comparisonListRepository repository.ComparisonListRepository
	calcDiscountService      services.CalcDiscountService
	hasher                   hasher.Hasher
	logger                   *logger.Logger
}

func NewUserServiceImplementation(userRepository repository.UserRepository, comparisonListRepository repository.ComparisonListRepository, calcDiscountService services.CalcDiscountService, logger *logger.Logger) services.UserService {
	return &userServiceImplementation{
		userRepository:           userRepository,
		comparisonListRepository: comparisonListRepository,
		calcDiscountService:      calcDiscountService,
		hasher:                   &implementation.BcryptHasher{},
		logger:                   logger,
	}
}

func (u *userServiceImplementation) Create(user *models.User, password string) error {
	fields := map[string]interface{}{"user_login": user.Login}
	_, err := u.userRepository.Get(user.Login)
	if err != nil && err != repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(fields).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	} else if err == nil {
		u.logger.WithFields(fields).Error(serviceErrors.UserCreateFailed.Error() + serviceErrors.UserAlreadyExists.Error())
		return serviceErrors.UserAlreadyExists
	}

	hashPassword, err := u.hasher.GetHash(password)
	if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}
	user.Password = string(hashPassword)

	err = u.userRepository.Create(user)
	if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}

	newUser, err := u.userRepository.Get(user.Login)
	if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}
	err = u.comparisonListRepository.Create(&models.ComparisonList{UserId: newUser.UserId})
	if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return err
	}

	u.logger.WithFields(fields).Info("user create completed")

	return nil
}

func (u *userServiceImplementation) Get(login string, password string) (*models.User, error) {
	fields := map[string]interface{}{"user_login": login}
	user, err := u.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(fields).Error(serviceErrors.UserGetFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.UserGetFailed.Error() + err.Error())
		return nil, err
	}

	if !u.hasher.Check(user.Password, password) {
		u.logger.WithFields(fields).Error(serviceErrors.UserGetFailed.Error() + serviceErrors.InvalidPassword.Error())
		return nil, serviceErrors.InvalidPassword
	}

	u.logger.WithFields(fields).Info("user get completed")

	return user, nil
}

func (u *userServiceImplementation) GetComparisonList(id uint64) (*models.ComparisonList, []models.Instrument, error) {
	fields := map[string]interface{}{"user_id": id}
	user, err := u.userRepository.GetById(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return nil, nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	fields = map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}
	comparisonList, err := u.comparisonListRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return nil, nil, serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	fields = map[string]interface{}{"user_login": user.Login, "user_id": user.UserId, "comparisonList_id": comparisonList.ComparisonListId}
	instruments, err := u.comparisonListRepository.GetInstruments(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, nil, nil
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	fields = map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}
	instruments, err = u.calcDiscountService.CalcDiscount(user, instruments)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, instruments, nil
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	var totalPrice uint64
	for i := range instruments {
		totalPrice += instruments[i].Price
	}

	fieldsToUpdate := models.ComparisonListFieldsToUpdate{}

	fieldsToUpdate[models.ComparisonListFieldTotalPrice] = totalPrice
	fieldsToUpdate[models.ComparisonListFieldAmount] = len(instruments)

	err = u.comparisonListRepository.Update(comparisonList.ComparisonListId, fieldsToUpdate)
	if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListUpdateFailed.Error() + err.Error())
		return nil, nil, err
	}

	comparisonList, err = u.comparisonListRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListGetFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return nil, nil, serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		u.logger.WithFields(fields).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListGetFailed.Error() + err.Error())
		return nil, nil, err
	}

	fields = map[string]interface{}{"user_id": id}
	u.logger.WithFields(fields).Info("user get comparisonList completed")

	return comparisonList, instruments, nil
}
