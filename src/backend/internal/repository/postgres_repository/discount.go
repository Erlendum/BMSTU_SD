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

type DiscountPostgres struct {
	DiscountId   uint64    `db:"discount_id"`
	InstrumentId uint64    `db:"instrument_id"`
	UserId       uint64    `db:"user_id"`
	Amount       uint64    `db:"discount_amount"`
	Type         string    `db:"discount_type"`
	DateBegin    time.Time `db:"discount_date_begin"`
	DateEnd      time.Time `db:"discount_date_end"`
}

type DiscountPostgresRepository struct {
	db *sqlx.DB
}

func NewDiscountPostgresRepository(db *sqlx.DB) repository.DiscountRepository {
	return &DiscountPostgresRepository{db: db}
}

func (i *DiscountPostgresRepository) Create(discount *models.Discount) error {
	query := `insert into store.discounts (discount_id, instrument_id, user_id, discount_amount,
											 discount_type, discount_date_begin, discount_date_end) values
											 ($1, $2, $3, $4, $5, $6, $7);`
	_, err := i.db.Exec(query, discount.DiscountId, discount.InstrumentId, discount.UserId,
		discount.Amount, discount.Type, discount.DateBegin, discount.DateEnd)
	if err != nil {
		return err
	}
	return nil
}

func (i *DiscountPostgresRepository) discountFieldToDBField(field models.DiscountField) string {
	switch field {
	case models.DiscountFieldInstrumentId:
		return "instrument_id"
	case models.DiscountFieldUserId:
		return "user_id"
	case models.DiscountFieldAmount:
		return "discount_amount"
	case models.DiscountFieldType:
		return "discount_type"
	case models.DiscountFieldDateBegin:
		return "discount_date_begin"
	case models.DiscountFieldDateEnd:
		return "discount_date_end"
	}
	return ""
}

func (i *DiscountPostgresRepository) Update(id uint64, fieldsToUpdate models.DiscountFieldsToUpdate) error {
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		updateFields[i.discountFieldToDBField(key)] = value
	}

	query, fields := queries.CreateUpdateQuery("store.discounts", updateFields)

	fields = append(fields, id)
	query += ` where discount_id = $` + strconv.Itoa(len(fields)) + ";"

	res, err := i.db.Exec(query, fields...)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (i *DiscountPostgresRepository) Delete(id uint64) error {
	query := `delete from store.discounts where discount_id = $1`
	res, err := i.db.Exec(query, id)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (i *DiscountPostgresRepository) Get(id uint64) (*models.Discount, error) {
	query := `select * from store.discounts where discount_id = $1`
	discountPostgres := &DiscountPostgres{}

	err := i.db.Get(discountPostgres, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	discount := &models.Discount{}
	err = copier.Copy(discount, discountPostgres)
	if err != nil {
		return nil, err
	}

	return discount, nil
}

func (i *DiscountPostgresRepository) GetList() ([]models.Discount, error) {
	query := `select * from store.discounts;`

	var discountsPostgres []DiscountPostgres
	var discounts []models.Discount
	err := i.db.Select(&discountsPostgres, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range discountsPostgres {
		discount := &models.Discount{}
		err = copier.Copy(discount, &discountsPostgres[i])
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, *discount)
	}
	return discounts, nil
}

func (i *DiscountPostgresRepository) GetSpecificList(instrumentId uint64, userId uint64) ([]models.Discount, error) {
	query := `select * from store.discounts where instrument_id = $1 and user_id = $2;`

	var discountsPostgres []DiscountPostgres
	var discounts []models.Discount
	err := i.db.Select(&discountsPostgres, query, instrumentId, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range discountsPostgres {
		discount := &models.Discount{}
		err = copier.Copy(discount, &discountsPostgres[i])
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, *discount)
	}
	return discounts, nil
}
