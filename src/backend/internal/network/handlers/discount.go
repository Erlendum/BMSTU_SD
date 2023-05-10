package handlers

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/handlerErrors"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var (
	DEFAULTTIME = time.Date(2006, 01, 02, 15, 04, 05, 0, time.UTC)
)

type DiscountHandler struct {
	service services.DiscountService
}

func NewDiscountHandler(service services.DiscountService) *DiscountHandler {
	return &DiscountHandler{service: service}
}

func (h *DiscountHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	var discount models.Discount
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	err := h.service.Create(&discount, login)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *DiscountHandler) CreateForAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	var discount models.Discount
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	err := h.service.CreateForAll(&discount, login)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *DiscountHandler) GetList(w http.ResponseWriter, r *http.Request) {
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
	e, err := h.service.GetList()
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	strLimit := r.URL.Query().Get("limit")
	limit := len(e)
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < -1 {
			http.Error(w, "limit query parameter is no valid number", http.StatusBadRequest)
			return
		}
	}

	strOffset := r.URL.Query().Get("offset")
	skip := 0
	if strOffset != "" {
		skip, err = strconv.Atoi(strOffset)
		if err != nil || skip < -1 {
			http.Error(w, "offset query parameter is no valid number", http.StatusBadRequest)
			return
		}
	}
	structure := make(map[string]any)
	structure["discounts"] = e[skip : skip+limit]
	structure["limit"] = limit
	structure["skip"] = skip
	sendMapResponse(w, http.StatusOK, structure)
}

func (h *DiscountHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	e := 0
	strId := r.URL.Query().Get("id")
	login := r.URL.Query().Get("login")
	id, err := strconv.Atoi(strId)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}
	err = h.service.Delete(uint64(id), login)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *DiscountHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	var discount models.Discount
	fields := make(models.DiscountFieldsToUpdate)
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	if discount.InstrumentId != 0 {
		fields[models.DiscountFieldInstrumentId] = discount.InstrumentId
	}
	if discount.UserId != 0 {
		fields[models.DiscountFieldUserId] = discount.UserId
	}
	if discount.Amount != 0 {
		fields[models.DiscountFieldAmount] = discount.Amount
	}

	if discount.DateBegin != DEFAULTTIME {
		fields[models.DiscountFieldDateBegin] = discount.DateBegin
	}

	if discount.DateEnd != DEFAULTTIME {
		fields[models.DiscountFieldDateBegin] = discount.DateEnd
	}

	err := h.service.Update(discount.DiscountId, login, fields)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}
