package handlers

import (
	"backend/internal/services"
	"log"
)

type ComparisonListHandler struct {
	service services.ComparisonListService
}

func NewComparisonListHandler(service services.ComparisonListService) *ComparisonListHandler {
	return &ComparisonListHandler{service: service}
}

func (h *ComparisonListHandler) AddInstrument(id uint64, instrumentId uint64) string {

	e := 0
	err := h.service.AddInstrument(id, instrumentId)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *ComparisonListHandler) DeleteInstrument(id uint64, instrumentId uint64) string {
	e := 0
	err := h.service.DeleteInstrument(id, instrumentId)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}
