package handlers

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/handlerErrors"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	var user models.User
	if r.Method != "POST" {
		sendErrorResponse(w, &ErrorModel{
			Error:          handlerErrors.PostExpectedError.Error(),
			HTTPStatusCode: http.StatusMethodNotAllowed,
		})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		})
		return
	}
	e := 0
	err := h.service.Create(&user, user.Password)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	sendResponse(w, http.StatusOK, e)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	var user *models.User
	login := r.URL.Query().Get("login")
	password := r.URL.Query().Get("password")
	user, err := h.service.Get(login, password)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	structure := make(map[string]any)
	structure["user"] = user
	sendMapResponse(w, http.StatusOK, structure)
}

func (h *UserHandler) GetComparisonList(w http.ResponseWriter, r *http.Request) {
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
	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}
	comparisonList, instruments, err := h.service.GetComparisonList(uint64(id))
	if err != nil {
		sendErrorResponse(w, &ErrorModel{
			Error:          err.Error(),
			HTTPStatusCode: http.StatusServiceUnavailable,
		})
		return
	}

	structure := make(map[string]any)
	structure["comparisonList"] = comparisonList
	structure["cartInstruments"] = instruments
	sendMapResponse(w, http.StatusOK, structure)
}
