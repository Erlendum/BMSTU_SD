package handlers

import (
	"backend/internal/pkg/errors/handlerErrors"
	"backend/internal/services"
	"net/http"
	"strconv"
)

type ComparisonListHandler struct {
	service services.ComparisonListService
}

func NewComparisonListHandler(service services.ComparisonListService) *ComparisonListHandler {
	return &ComparisonListHandler{service: service}
}

func (h *ComparisonListHandler) AddInstrument(w http.ResponseWriter, r *http.Request) {
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
	strInstrumentId := r.URL.Query().Get("instrumentId")
	id, err := strconv.Atoi(strId)

	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}

	instrumentId, err := strconv.Atoi(strInstrumentId)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}

	err = h.service.AddInstrument(uint64(id), uint64(instrumentId))
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *ComparisonListHandler) DeleteInstrument(w http.ResponseWriter, r *http.Request) {
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
	strInstrumentId := r.URL.Query().Get("instrumentId")
	id, err := strconv.Atoi(strId)

	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}

	instrumentId, err := strconv.Atoi(strInstrumentId)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
	}

	err = h.service.DeleteInstrument(uint64(id), uint64(instrumentId))
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}
