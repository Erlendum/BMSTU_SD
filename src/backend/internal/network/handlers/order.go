package handlers

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/handlerErrors"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	var e uint64
	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)

	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}

	e, err = h.service.Create(uint64(id))
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, int(e))
}

func (h *OrderHandler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method != "GET" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.GetExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	strId := r.URL.Query().Get("user_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}
	e, err := h.service.GetList(uint64(id))
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	structure := make(map[string]any)
	structure["orders"] = e
	sendMapResponse(w, http.StatusOK, structure)
}

func (h *OrderHandler) GetListForAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method != "GET" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.GetExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}

	e, err := h.service.GetListForAll()
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	structure := make(map[string]any)
	structure["orders"] = e
	sendMapResponse(w, http.StatusOK, structure)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	var order models.Order
	fields := make(models.OrderFieldsToUpdate)
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	if order.UserId != 0 {
		fields[models.OrderFieldUserId] = order.UserId
	}
	if order.Price != 0 {
		fields[models.OrderFieldPrice] = order.Price
	}
	if len(order.Status) != 0 {
		fields[models.OrderFieldStatus] = order.Status
	}

	if order.Time != DEFAULTTIME {
		fields[models.OrderFieldTime] = order.Time
	}

	if order.UserId != 0 {
		fields[models.OrderFieldUserId] = order.UserId
	}

	err := h.service.Update(order.OrderId, login, fields)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}
