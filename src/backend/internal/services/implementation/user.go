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
	u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Info("user create called")

	_, err := u.userRepository.Get(user.Login)
	if err != nil && err != repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	} else if err == nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.UserCreateFailed.Error() + serviceErrors.UserAlreadyExists.Error())
		return serviceErrors.UserAlreadyExists
	}

	hashPassword, err := u.hasher.GetHash(password)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}
	user.Password = string(hashPassword)

	err = u.userRepository.Create(user)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}

	newUser, err := u.userRepository.Get(user.Login)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.UserCreateFailed.Error() + err.Error())
		return err
	}
	err = u.comparisonListRepository.Create(&models.ComparisonList{UserId: newUser.UserId})
	if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return err
	}
	return nil
}

func (u *userServiceImplementation) Get(login string, password string) (*models.User, error) {
	u.logger.WithFields(map[string]interface{}{"user_login": login}).Info("user get called")

	user, err := u.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(map[string]interface{}{"user_login": login}).Error(serviceErrors.UserGetFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": login}).Error(serviceErrors.UserGetFailed.Error() + err.Error())
		return nil, err
	}

	if !u.hasher.Check(user.Password, password) {
		u.logger.WithFields(map[string]interface{}{"user_login": login}).Error(serviceErrors.UserGetFailed.Error() + serviceErrors.InvalidPassword.Error())
		return nil, serviceErrors.InvalidPassword
	}

	return user, nil
}

func (u *userServiceImplementation) GetComparisonList(id uint64) (*models.ComparisonList, []models.Instrument, error) {
	u.logger.WithFields(map[string]interface{}{"user_id": id}).Info("user get comparisonList called")

	user, err := u.userRepository.GetById(id)

	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return nil, nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login}).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	comparisonList, err := u.comparisonListRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return nil, nil, serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	instruments, err := u.comparisonListRepository.GetInstruments(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, nil, nil
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId, "comparisonList_id": comparisonList.ComparisonListId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	instruments, err = u.calcDiscountService.CalcDiscount(user, instruments)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, instruments, nil
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + err.Error())
		return nil, nil, err
	}

	var totalPrice uint64
	for i := range instruments {
		totalPrice += instruments[i].Price
	}

	fields := models.ComparisonListFieldsToUpdate{}

	fields[models.ComparisonListFieldTotalPrice] = totalPrice
	fields[models.ComparisonListFieldAmount] = len(instruments)

	err = u.comparisonListRepository.Update(comparisonList.ComparisonListId, fields)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListUpdateFailed.Error() + err.Error())
		return nil, nil, err
	}

	comparisonList, err = u.comparisonListRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListGetFailed.Error() + serviceErrors.ComparisonListDoesNotExists.Error())
		return nil, nil, serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		u.logger.WithFields(map[string]interface{}{"user_login": user.Login, "user_id": user.UserId}).Error(serviceErrors.ComparisonListCreateFailed.Error() + serviceErrors.ComparisonListGetFailed.Error() + err.Error())
		return nil, nil, err
	}

	return comparisonList, instruments, nil
}
