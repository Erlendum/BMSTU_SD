package CLIhandlers

import (
	"backend/internal/services"
	"log"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(userId uint64) (uint64, string) {
	e := 0
	res, err := h.service.Create(userId)
	if err != nil {
		log.Println(err)
		return 0, ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return res, Response(e)
}
