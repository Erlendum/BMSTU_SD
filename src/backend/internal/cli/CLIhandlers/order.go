package CLIhandlers

import (
	"backend/internal/models"
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

func (h *OrderHandler) GetList(userId uint64) string {

	e, err := h.service.GetList(userId)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}
	structure := make(map[string]any)
	structure["orders"] = e

	return MapResponse(structure)
}

func (h *OrderHandler) GetListForAll() string {

	e, err := h.service.GetListForAll()
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}
	structure := make(map[string]any)
	structure["orders"] = e

	return MapResponse(structure)
}

func (h *OrderHandler) Update(id uint64, login string, fieldsToUpdate models.OrderFieldsToUpdate) string {
	e := 0
	err := h.service.Update(id, login, fieldsToUpdate)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *OrderHandler) GetOrderElements(orderId uint64) string {

	e, err := h.service.GetOrderElements(orderId)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}
	structure := make(map[string]any)
	structure["order_elements"] = e

	return MapResponse(structure)
}
