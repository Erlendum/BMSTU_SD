package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(user models.User) string {
	e := 0
	err := h.service.Create(&user, user.Password)
	if err != nil {
		log.Println(err)
		return ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}

	return Response(e)
}

func (h *UserHandler) Get(login, password string) (*models.User, string) {
	e := 0
	user, err := h.service.Get(login, password)
	if err != nil {
		log.Println(err)
		return nil, ErrorResponse(&ErrorModel{
			Error: err.Error(),
		})
	}
	return user, Response(e)
}
