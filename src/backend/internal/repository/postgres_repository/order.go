package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/repository"
	"github.com/jmoiron/sqlx"
	"time"
)

type OrderPostgres struct {
	OrderId uint64    `db:"order_id"`
	UserId  uint64    `db:"user_id"`
	Price   uint64    `db:"order_price"`
	Time    time.Time `db:"order_time"`
	Status  string    `db:"order_status"`
}

type OrderElementPostgres struct {
	OrderElementId uint64 `db:"order_element_id"`
	InstrumentId   uint64 `db:"instrument_id"`
	OrderId        uint64 `db:"order_id"`
	Amount         uint64 `db:"order_element_amount"`
	Price          uint64 `db:"order_element_price"`
}

type OrderPostgresRepository struct {
	db *sqlx.DB
}

func NewOrderPostgresRepository(db *sqlx.DB) repository.OrderRepository {
	return &OrderPostgresRepository{db: db}
}

func (i *OrderPostgresRepository) Create(order *models.Order) (uint64, error) {
	query := `insert into store.orders (order_time, order_price, user_id) values
											 ($1, $2, $3) returning order_id;`
	var id uint64
	err := i.db.Get(&id, query, order.Time, order.Price, order.UserId)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (i *OrderPostgresRepository) CreateOrderElement(element *models.OrderElement) error {
	query := `insert into store.order_elements (order_element_amount, order_element_price, instrument_id, order_id) values
											 ($1, $2, $3, $4);`
	_, err := i.db.Exec(query, element.Amount, element.Price, element.InstrumentId, element.OrderId)
	if err != nil {
		return err
	}

	return nil
}
