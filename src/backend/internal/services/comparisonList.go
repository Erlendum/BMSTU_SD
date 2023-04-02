package services

type ComparisonListService interface {
	AddInstrument(id uint64, instrumentId uint64) error
	DeleteInstrument(id uint64, instrumentId uint64) error
}
