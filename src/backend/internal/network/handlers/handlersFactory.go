package handlers

import "backend/internal/services"

type HandlersServicesFields struct {
	InstrumentService     services.InstrumentService
	UserService           services.UserService
	ComparisonListService services.ComparisonListService
	DiscountService       services.DiscountService
	OrderService          services.OrderService
}

type Handlers struct {
	InstrumentHandler     *InstrumentHandler
	UserHandler           *UserHandler
	ComparisonListHandler *ComparisonListHandler
	DiscountHandler       *DiscountHandler
	OrderHandler          *OrderHandler
}

func NewHandlers(services HandlersServicesFields) *Handlers {
	return &Handlers{
		InstrumentHandler:     NewInstrumentHandler(services.InstrumentService),
		UserHandler:           NewUserHandler(services.UserService),
		ComparisonListHandler: NewComparisonListHandler(services.ComparisonListService),
		DiscountHandler:       NewDiscountHandler(services.DiscountService),
		OrderHandler:          NewOrderHandler(services.OrderService),
	}
}
