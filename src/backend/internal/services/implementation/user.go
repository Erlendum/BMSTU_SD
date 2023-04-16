package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/hasher"
	"backend/internal/pkg/hasher/implementation"
	"backend/internal/repository"
	"backend/internal/services"
)

type userServiceImplementation struct {
	userRepository           repository.UserRepository
	comparisonListRepository repository.ComparisonListRepository
	calcDiscountService      services.CalcDiscountService
	hasher                   hasher.Hasher
}

func NewUserServiceImplementation(userRepository repository.UserRepository, comparisonListRepository repository.ComparisonListRepository, calcDiscountService services.CalcDiscountService) services.UserService {
	return &userServiceImplementation{
		userRepository:           userRepository,
		comparisonListRepository: comparisonListRepository,
		calcDiscountService:      calcDiscountService,
		hasher:                   &implementation.BcryptHasher{},
	}
}

func (u *userServiceImplementation) Create(user *models.User, password string) error {
	_, err := u.userRepository.Get(user.Login)
	if err != nil && err != repositoryErrors.ObjectDoesNotExists {
		return err
	} else if err == nil {
		return serviceErrors.UserAlreadyExists
	}

	hashPassword, err := u.hasher.GetHash(password)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	err = u.userRepository.Create(user)
	if err != nil {
		return err
	}

	newUser, err := u.userRepository.Get(user.Login)
	if err != nil {
		return err
	}

	return u.comparisonListRepository.Create(&models.ComparisonList{UserId: newUser.UserId})
}

func (u *userServiceImplementation) Get(login string, password string) (*models.User, error) {
	user, err := u.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		return nil, err
	}

	if !u.hasher.Check(user.Password, password) {
		return nil, serviceErrors.InvalidPassword
	}

	return user, nil
}

func (u *userServiceImplementation) GetComparisonList(id uint64) (*models.ComparisonList, []models.Instrument, error) {
	user, err := u.userRepository.GetById(id)

	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, nil, serviceErrors.UserDoesNotExists
	} else if err != nil {
		return nil, nil, err
	}

	comparisonList, err := u.comparisonListRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, nil, serviceErrors.ComparisonListDoesNotExists
	} else if err != nil {
		return nil, nil, err
	}

	instruments, err := u.comparisonListRepository.GetInstruments(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, nil, nil
	} else if err != nil {
		return nil, nil, err
	}

	instruments, err = u.calcDiscountService.CalcDiscount(user, instruments)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return comparisonList, instruments, nil
	} else if err != nil {
		return nil, nil, err
	}

	return comparisonList, instruments, nil
}
