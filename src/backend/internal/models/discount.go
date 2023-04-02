package models

import "time"

type DiscountField int
type DiscountFieldsToUpdate map[DiscountField]any

const (
	DiscountFieldInstrumentId = DiscountField(iota)
	DiscountFieldUserId
	DiscountFieldAmount
	DiscountFieldType
	DiscountFieldDateBegin
	DiscountFieldDateEnd
)

type Discount struct {
	DiscountId   uint64
	InstrumentId uint64
	UserId       uint64
	Amount       uint64
	Type         string
	DateBegin    time.Time
	DateEnd      time.Time
}
