package CLIhandlers

import (
	"backend/internal/models"
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
	structure["discounts"] = e[start:stop]
	structure["limit"] = limit
	structure["skip"] = skip
	return MapResponse(structure)
}

func (h *DiscountHandler) Create(discount models.Discount, login string) string {
	e := 0
	err := h.service.Create(&discount, login)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *DiscountHandler) CreateForAll(discount models.Discount, login string) string {
	e := 0
	err := h.service.CreateForAll(&discount, login)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *DiscountHandler) Delete(id uint64, login string) string {
	e := 0
	err := h.service.Delete(id, login)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *DiscountHandler) Update(id uint64, login string, fieldsToUpdate models.DiscountFieldsToUpdate) string {
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
