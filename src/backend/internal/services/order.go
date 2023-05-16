package services

type OrderService interface {
	Create(userId uint64) (uint64, error)
}
