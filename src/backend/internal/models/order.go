package models

import "time"

type (
	OrderField          int
	OrderFieldsToUpdate map[OrderField]any
)

const (
	OrderFieldUserId = OrderField(iota)
	OrderFieldPrice
	OrderFieldTime
	OrderFieldStatus
)

type Order struct {
	OrderId uint64
	UserId  uint64
	Price   uint64
	Time    time.Time
	Status  string
}
