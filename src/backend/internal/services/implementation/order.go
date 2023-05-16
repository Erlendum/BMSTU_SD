package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"time"
)

type orderServiceImplementation struct {
	orderRepository          repository.OrderRepository
	comparisonListRepository repository.ComparisonListRepository
}

func NewOrderServiceImplementation(orderRepository repository.OrderRepository, comparisonListRepository repository.ComparisonListRepository) services.OrderService {
	return &orderServiceImplementation{
		orderRepository:          orderRepository,
		comparisonListRepository: comparisonListRepository,
	}
}

func (o *orderServiceImplementation) Create(userId uint64) (uint64, error) {
	comparisonList, err := o.comparisonListRepository.Get(userId)
	if err != nil {
		return 0, err
	}
	instruments, err := o.comparisonListRepository.GetInstruments(userId)
	if err != nil {
		return 0, err
	}

	var totalPrice uint64
	for i := range instruments {
		totalPrice += instruments[i].Price
	}

	order := models.Order{UserId: userId, Price: totalPrice, Time: time.Now()}
	orderId, err := o.orderRepository.Create(&order)
	if err != nil {
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
			return 0, err
		}
	}

	err = o.comparisonListRepository.Clear(comparisonList.ComparisonListId)
	if err != nil {
		return 0, err
	}

	fields := make(models.ComparisonListFieldsToUpdate)
	fields[models.ComparisonListFieldTotalPrice] = 0
	fields[models.ComparisonListFieldAmount] = 0
	err = o.comparisonListRepository.Update(comparisonList.ComparisonListId, fields)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}
