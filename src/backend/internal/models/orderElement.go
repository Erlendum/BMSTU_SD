package models

type (
	OrderElementField          int
	OrderElementFieldsToUpdate map[OrderElementField]any
)

const (
	OrderElementFieldInstrumentId = OrderElementField(iota)
	OrderElementFieldOrderId
	OrderElementFieldAmount
	OrderElementFieldPrice
)

type OrderElement struct {
	OrderElementId uint64
	InstrumentId   uint64
	OrderId        uint64
	Amount         uint64
	Price          uint64
}
