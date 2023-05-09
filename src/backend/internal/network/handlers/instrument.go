package handlers

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/handlerErrors"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type InstrumentHandler struct {
	service services.InstrumentService
}

func NewInstrumentHandler(service services.InstrumentService) *InstrumentHandler {
	return &InstrumentHandler{service: service}
}

func (h *InstrumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	var instrument models.Instrument
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&instrument); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	err := h.service.Create(&instrument, login)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *InstrumentHandler) GetList(w http.ResponseWriter, r *http.Request) {
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
	structure["instruments"] = e[skip : skip+limit]
	structure["limit"] = limit
	structure["skip"] = skip
	sendMapResponse(w, http.StatusOK, structure)
}

func (h *InstrumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *InstrumentHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	var instrument models.Instrument
	fields := make(models.InstrumentFieldsToUpdate)
	if err := json.NewDecoder(r.Body).Decode(&instrument); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	login := r.URL.Query().Get("login")
	if len(instrument.Brand) != 0 {
		fields[models.InstrumentFieldBrand] = instrument.Brand
	}
	if len(instrument.Name) != 0 {
		fields[models.InstrumentFieldName] = instrument.Name
	}
	if instrument.Price != 0 {
		fields[models.InstrumentFieldPrice] = instrument.Price
	}

	if len(instrument.Material) != 0 {
		fields[models.InstrumentFieldMaterial] = instrument.Material
	}

	if len(instrument.Type) != 0 {
		fields[models.InstrumentFieldType] = instrument.Type
	}

	if len(instrument.Img) != 0 {
		fields[models.InstrumentFieldImg] = instrument.Img
	}

	err := h.service.Update(instrument.InstrumentId, login, fields)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}
