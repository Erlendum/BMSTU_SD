package models

import "time"

type UserField int
type UserFieldsToUpdate map[UserField]any

const (
	UserFieldLogin = UserField(iota)
	UserFieldFio
	UserFieldDateBirth
	UserFieldGender
	IsAdmin
)

type UserGender string

const (
	MaleUserGender   = UserGender("Male")
	FemaleUserGender = UserGender("Female")
)

type User struct {
	UserId    uint64
	Login     string
	Password  string
	Fio       string
	DateBirth time.Time
	Gender    UserGender
	IsAdmin   bool
}
