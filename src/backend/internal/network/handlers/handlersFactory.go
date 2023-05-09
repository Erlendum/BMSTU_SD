package handlers

import "backend/internal/services"

type HandlersServicesFields struct {
	InstrumentService services.InstrumentService
	UserService       services.UserService
}

type Handlers struct {
	InstrumentHandler *InstrumentHandler
	UserHandler       *UserHandler
}

func NewHandlers(services HandlersServicesFields) *Handlers {
	return &Handlers{
		InstrumentHandler: NewInstrumentHandler(services.InstrumentService),
		UserHandler:       NewUserHandler(services.UserService),
	}
}
