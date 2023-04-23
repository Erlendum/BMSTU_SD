package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/services"
)

type discountServiceImplementation struct {
	discountRepository repository.DiscountRepository
	userRepository     repository.UserRepository
	logger             *logger.Logger
}

func NewDiscountServiceImplementation(discountRepository repository.DiscountRepository, userRepository repository.UserRepository, logger *logger.Logger) services.DiscountService {
	return &discountServiceImplementation{
		discountRepository: discountRepository,
		userRepository:     userRepository,
		logger:             logger,
	}
}

func (d *discountServiceImplementation) Create(discount *models.Discount, login string) error {
	d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Info("discount create called")
	user, err := d.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountCreateFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return serviceErrors.UserDoesNotExists
	} else if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountCreateFailed.Error() + err.Error())
		return err
	}

	if !user.IsAdmin {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountCreateFailed.Error() + serviceErrors.UserCanNotCreateDiscount.Error())
		return serviceErrors.UserCanNotCreateDiscount
	}

	err = d.discountRepository.Create(discount)
	if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountCreateFailed.Error() + err.Error())
		return err
	}
	return nil
}

func (d *discountServiceImplementation) CreateForAll(discount *models.Discount, login string) error {
	d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Info("discount create for all called")

	user, err := d.userRepository.Get(login)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountForAllCreateFailed.Error() + serviceErrors.UserDoesNotExists.Error())
		return serviceErrors.UserDoesNotExists
	} else if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountForAllCreateFailed.Error() + err.Error())
		return err
	}

	if !user.IsAdmin {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountForAllCreateFailed.Error() + serviceErrors.UserCanNotCreateDiscount.Error())
		return serviceErrors.UserCanNotCreateDiscount
	}

	users, err := d.userRepository.GetList()
	if err != nil {
		return err
	}
	for i := range users {
		userId := users[i].UserId
		discount.UserId = userId
		err = d.discountRepository.Create(discount)
		if err != nil {
			d.logger.WithFields(map[string]interface{}{"user_login": login, "user_id": discount.UserId, "discount_type": discount.Type, "discount_dateBegin": discount.DateBegin, "discount_dateEnd": discount.DateEnd}).Error(serviceErrors.DiscountForAllCreateFailed.Error() + err.Error())
			return err
		}
	}

	return nil
}

func (d *discountServiceImplementation) Update(id uint64, login string, fieldsToUpdate models.DiscountFieldsToUpdate) error {
	d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Info("discount update called")

	canUpdate, err := d.checkCanUserChangeDiscount(id, login)
	if err != nil {
		return err
	} else if !canUpdate {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Error(serviceErrors.DiscountUpdateFailed.Error() + serviceErrors.UserCanNotUpdateDiscount.Error())
		return serviceErrors.UserCanNotUpdateDiscount
	}

	err = d.discountRepository.Update(id, fieldsToUpdate)
	if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Error(serviceErrors.DiscountUpdateFailed.Error() + err.Error())
		return err
	}
	return nil
}

func (d *discountServiceImplementation) Delete(id uint64, login string) error {
	d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Info("discount delete called")

	canDelete, err := d.checkCanUserChangeDiscount(id, login)
	if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Error(serviceErrors.DiscountDeleteFailed.Error() + err.Error())
		return err
	} else if !canDelete {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Error(serviceErrors.DiscountDeleteFailed.Error() + serviceErrors.UserCanNotDeleteDiscount.Error())
		return serviceErrors.UserCanNotDeleteDiscount
	}

	err = d.discountRepository.Delete(id)
	if err != nil {
		d.logger.WithFields(map[string]interface{}{"user_login": login, "discount_id": id}).Error(serviceErrors.DiscountDeleteFailed.Error() + err.Error())
	}
	return nil
}

func (d *discountServiceImplementation) Get(id uint64) (*models.Discount, error) {
	d.logger.WithFields(map[string]interface{}{"discount_id": id}).Info("discount get called")

	discount, err := d.discountRepository.Get(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		d.logger.WithFields(map[string]interface{}{"discount_id": id}).Error(serviceErrors.DiscountGetFailed.Error() + serviceErrors.DiscountDoesNotExists.Error())
		return nil, serviceErrors.DiscountDoesNotExists
	} else if err != nil {
		d.logger.WithFields(map[string]interface{}{"discount_id": id}).Error(serviceErrors.DiscountGetFailed.Error() + err.Error())
		return nil, err
	}

	return discount, nil
}

func (d *discountServiceImplementation) GetList() ([]models.Discount, error) {
	d.logger.Info("discount get list called")
	discounts, err := d.discountRepository.GetList()
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		d.logger.Error(serviceErrors.DiscountsListGetFailed.Error() + serviceErrors.DiscountsDoesNotExists.Error())
		return nil, serviceErrors.DiscountsDoesNotExists
	} else if err != nil {
		d.logger.Error(serviceErrors.DiscountsListGetFailed.Error() + err.Error())
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
