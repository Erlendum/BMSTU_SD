package handlers

import (
	"backend/internal/services"
	"log"
)

type DiscountHandler struct {
	service services.DiscountService
}

func NewDiscountHandler(service services.DiscountService) *DiscountHandler {
	return &DiscountHandler{service: service}
}

func (h *DiscountHandler) GetList(skip, limit int) string {

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
