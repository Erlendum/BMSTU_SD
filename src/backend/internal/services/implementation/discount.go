package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/repository"
	"backend/internal/services"
)

type discountServiceImplementation struct {
	discountRepository repository.DiscountRepository
	userRepository     repository.UserRepository
}

func NewDiscountServiceImplementation(discountRepository repository.DiscountRepository, userRepository repository.UserRepository) services.DiscountService {
	return &discountServiceImplementation{
		discountRepository: discountRepository,
		userRepository:     userRepository,
	}
}

func (d *discountServiceImplementation) Create(discount *models.Discount, login string) error {
	user, err := d.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return serviceErrors.UserDoesNotExists
	} else if err != nil {
		return err
	}

	if !user.IsAdmin {
		return serviceErrors.UserCanNotCreateDiscount
	}

	if discount.UserId == 0 {
		users, err := d.userRepository.GetList()
		if err != nil {
			return err
		}
		for i := range users {
			userId := users[i].UserId
			discount.UserId = userId
			err = d.discountRepository.Create(discount)
			if err != nil {
				return nil
			}
		}
	} else {
		return d.discountRepository.Create(discount)
	}
	return nil
}

func (d *discountServiceImplementation) Update(id uint64, login string, fieldsToUpdate models.DiscountFieldsToUpdate) error {
	canUpdate, err := d.checkCanUserChangeDiscount(id, login)
	if err != nil {
		return err
	} else if !canUpdate {
		return serviceErrors.UserCanNotUpdateDiscount
	}

	return d.discountRepository.Update(id, fieldsToUpdate)
}

func (d *discountServiceImplementation) Delete(id uint64, login string) error {
	canDelete, err := d.checkCanUserChangeDiscount(id, login)
	if err != nil {
		return err
	} else if !canDelete {
		return serviceErrors.UserCanNotDeleteDiscount
	}

	return d.discountRepository.Delete(id)
}

func (d *discountServiceImplementation) Get(id uint64) (*models.Discount, error) {
	discount, err := d.discountRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, serviceErrors.DiscountDoesNotExists
	} else if err != nil {
		return nil, err
	}

	return discount, nil
}

func (d *discountServiceImplementation) GetList() ([]models.Discount, error) {
	discounts, err := d.discountRepository.GetList()
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		return nil, serviceErrors.DiscountsDoesNotExists
	} else if err != nil {
		return nil, err
	}

	return discounts, nil
}

func (d *discountServiceImplementation) checkCanUserChangeDiscount(discountId uint64, login string) (bool, error) {
	_, err := d.discountRepository.Get(discountId)
	if err == repositoryErrors.ObjectDoesNotExists {
		return false, serviceErrors.DiscountDoesNotExists
	} else if err != nil {
		return false, err
	}

	user, err := d.userRepository.Get(login)
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
