package models

type InstrumentField int
type InstrumentFieldsToUpdate map[InstrumentField]any

const (
	InstrumentFieldName = InstrumentField(iota)
	InstrumentFieldPrice
	InstrumentFieldMaterial
	InstrumentFieldType
	InstrumentFieldBrand
)

type Instrument struct {
	InstrumentId uint64
	Name         string
	Price        uint64
	Material     string
	Type         string
	Brand        string
}
