package handlers

import "backend/internal/services"

type HandlersServicesFields struct {
	InstrumentService     services.InstrumentService
	UserService           services.UserService
	ComparisonListService services.ComparisonListService
}

type Handlers struct {
	InstrumentHandler     *InstrumentHandler
	UserHandler           *UserHandler
	ComparisonListHandler *ComparisonListHandler
}

func NewHandlers(services HandlersServicesFields) *Handlers {
	return &Handlers{
		InstrumentHandler:     NewInstrumentHandler(services.InstrumentService),
		UserHandler:           NewUserHandler(services.UserService),
		ComparisonListHandler: NewComparisonListHandler(services.ComparisonListService),
	}
}
