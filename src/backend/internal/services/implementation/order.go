package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/services"
	"time"
)

type orderServiceImplementation struct {
	orderRepository          repository.OrderRepository
	comparisonListRepository repository.ComparisonListRepository
	userRepository           repository.UserRepository
	logger                   *logger.Logger
}

func NewOrderServiceImplementation(orderRepository repository.OrderRepository, comparisonListRepository repository.ComparisonListRepository, userRepository repository.UserRepository, logger *logger.Logger) services.OrderService {
	return &orderServiceImplementation{
		orderRepository:          orderRepository,
		comparisonListRepository: comparisonListRepository,
		userRepository:           userRepository,
		logger:                   logger,
	}
}

func (o *orderServiceImplementation) Create(userId uint64) (uint64, error) {
	fields := map[string]interface{}{"user_id": userId}

	comparisonList, err := o.comparisonListRepository.Get(userId)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
		return 0, err
	}
	instruments, err := o.comparisonListRepository.GetInstruments(userId)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
		return 0, err
	}

	var totalPrice uint64
	for i := range instruments {
		totalPrice += instruments[i].Price
	}

	order := models.Order{UserId: userId, Price: totalPrice, Time: time.Now()}
	orderId, err := o.orderRepository.Create(&order)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
		return 0, err
	}

	orderElementsMap := make(map[uint64]*models.OrderElement)

	for i := range instruments {
		if orderElementsMap[instruments[i].InstrumentId] != nil {
			orderElementsMap[instruments[i].InstrumentId].Amount++
		} else {
			orderElement := models.OrderElement{InstrumentId: instruments[i].InstrumentId, OrderId: orderId, Price: instruments[i].Price, Amount: 1}
			orderElementsMap[instruments[i].InstrumentId] = &orderElement
		}
	}

	for _, value := range orderElementsMap {
		err = o.orderRepository.CreateOrderElement(value)
		if err != nil {
			o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
			return 0, err
		}
	}

	err = o.comparisonListRepository.Clear(comparisonList.ComparisonListId)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
		return 0, err
	}

	fieldsToUpdate := make(models.ComparisonListFieldsToUpdate)
	fieldsToUpdate[models.ComparisonListFieldTotalPrice] = 0
	fieldsToUpdate[models.ComparisonListFieldAmount] = 0
	err = o.comparisonListRepository.Update(comparisonList.ComparisonListId, fieldsToUpdate)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderCreateFailed.Error() + err.Error())
		return 0, err
	}

	o.logger.WithFields(fields).Info("order create completed")
	return orderId, nil
}

func (o *orderServiceImplementation) GetList(userId uint64) ([]models.Order, error) {
	fields := map[string]interface{}{"user_id": userId}

	orders, err := o.orderRepository.GetList(userId)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		o.logger.WithFields(fields).Error(serviceErrors.OrdersListGetFailed.Error() + serviceErrors.OrdersDoesNotExists.Error())
		return nil, serviceErrors.OrdersDoesNotExists
	} else if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrdersListGetFailed.Error() + err.Error())
		return nil, err
	}

	o.logger.WithFields(fields).Info("order get list completed")
	return orders, nil
}

func (o *orderServiceImplementation) GetListForAll() ([]models.Order, error) {
	orders, err := o.orderRepository.GetListForAll()
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		o.logger.Error(serviceErrors.OrdersListForAllGetFailed.Error() + serviceErrors.OrdersDoesNotExists.Error())
		return nil, serviceErrors.OrdersDoesNotExists
	} else if err != nil {
		o.logger.Error(serviceErrors.OrdersListForAllGetFailed.Error() + err.Error())
		return nil, err
	}

	o.logger.Info("order get list for all completed")
	return orders, nil
}

func (o *orderServiceImplementation) Update(id uint64, login string, fieldsToUpdate models.OrderFieldsToUpdate) error {
	fields := map[string]interface{}{"user_login": login, "order_id": id}
	canUpdate, err := o.checkCanUserChangeOrder(login)
	if err != nil {
		return err
	} else if !canUpdate {
		o.logger.WithFields(fields).Error(serviceErrors.OrderUpdateFailed.Error() + serviceErrors.UserCanNotUpdateOrder.Error())
		return serviceErrors.UserCanNotUpdateOrder
	}

	err = o.orderRepository.Update(id, fieldsToUpdate)
	if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrderUpdateFailed.Error() + err.Error())
		return err
	}
	o.logger.WithFields(fields).Info("order update completed")

	return nil
}

func (o *orderServiceImplementation) checkCanUserChangeOrder(login string) (bool, error) {
	user, err := o.userRepository.Get(login)
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

func (o *orderServiceImplementation) GetOrderElements(id uint64) ([]models.OrderElement, error) {
	fields := map[string]interface{}{"order_id": id}

	orders, err := o.orderRepository.GetOrderElements(id)
	if err != nil && err == repositoryErrors.ObjectDoesNotExists {
		o.logger.WithFields(fields).Error(serviceErrors.OrdersListGetFailed.Error() + serviceErrors.OrdersDoesNotExists.Error())
		return nil, serviceErrors.OrdersDoesNotExists
	} else if err != nil {
		o.logger.WithFields(fields).Error(serviceErrors.OrdersListGetFailed.Error() + err.Error())
		return nil, err
	}

	o.logger.WithFields(fields).Info("get order elements completed")
	return orders, nil
}
