package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
)

type InstrumentHandler struct {
	service services.InstrumentService
}

func NewInstrumentHandler(service services.InstrumentService) *InstrumentHandler {
	return &InstrumentHandler{service: service}
}

func (h *InstrumentHandler) Create(instrument models.Instrument) string {
	e := 0
	err := h.service.Create(&instrument, "admin")
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *InstrumentHandler) GetList(skip, limit int) string {

	e, err := h.service.GetList()
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}
	start := skip
	stop := skip + limit
	if skip+limit > len(e) {
		stop = len(e)
	}
	structure := make(map[string]any)
	structure["instruments"] = e[start:stop]
	structure["limit"] = limit
	structure["skip"] = skip
	return MapResponse(structure)
}
