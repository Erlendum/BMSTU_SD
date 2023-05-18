package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/queries"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"strconv"
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

func (i *OrderPostgresRepository) GetList(userId uint64) ([]models.Order, error) {
	query := `select * from store.orders where user_id = $1`

	var ordersPostges []OrderPostgres
	var orders []models.Order
	err := i.db.Select(&ordersPostges, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range ordersPostges {
		order := &models.Order{}
		err = copier.Copy(order, &ordersPostges[i])
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (i *OrderPostgresRepository) GetListForAll() ([]models.Order, error) {
	query := `select * from store.orders;`

	var ordersPostgres []OrderPostgres
	var orders []models.Order
	err := i.db.Select(&ordersPostgres, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range ordersPostgres {
		order := &models.Order{}
		err = copier.Copy(order, &ordersPostgres[i])
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (i *OrderPostgresRepository) orderFieldToDBField(field models.OrderField) (string, error) {
	switch field {
	case models.OrderFieldUserId:
		return "user_id", nil
	case models.OrderFieldPrice:
		return "order_price", nil
	case models.OrderFieldTime:
		return "order_time", nil
	case models.OrderFieldStatus:
		return "order_status", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *OrderPostgresRepository) Update(id uint64, fieldsToUpdate models.OrderFieldsToUpdate) error {
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.orderFieldToDBField(key)
		if err != nil {
			return err
		}
		updateFields[field] = value
	}

	query, fields := queries.CreateUpdateQuery("store.orders", updateFields)

	fields = append(fields, id)
	query += ` where order_id = $` + strconv.Itoa(len(fields)) + ";"

	res, err := i.db.Exec(query, fields...)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}
