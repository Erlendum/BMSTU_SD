package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

type ComparisonListPostgres struct {
	ComparisonListId uint64 `db:"comparisonlist_id"`
	UserId           uint64 `db:"user_id"`
	TotalPrice       uint64 `db:"comparisonlist_total_price"`
	Amount           uint64 `db:"comparisonlist_amount"`
}

type ComparisonListPostgresRepository struct {
	db *sqlx.DB
}

func NewComparisonListPostgresRepository(db *sqlx.DB) repository.ComparisonListRepository {
	return &ComparisonListPostgresRepository{db: db}
}

func (i *ComparisonListPostgresRepository) Create(comparisonList *models.ComparisonList) error {
	query := `insert into store.comparisonLists (comparisonList_id, user_id, comparisonList_total_price, comparisonList_amount) values
											 ($1, $2, $3, $4);`
	_, err := i.db.Exec(query, comparisonList.ComparisonListId, comparisonList.UserId, comparisonList.TotalPrice,
		comparisonList.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListPostgresRepository) AddInstrument(id uint64, instrumentId uint64) error {
	query := `insert into store.comparisonLists_instruments (comparisonList_id, instrument_id) values
											 ($1, $2);`
	_, err := i.db.Exec(query, id, instrumentId)
	if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListPostgresRepository) DeleteInstrument(id uint64, instrumentId uint64) error {
	query := `delete from store.comparisonLists_instruments where comparisonList_id = $1 and instrument_id = $2;`
	res, err := i.db.Exec(query, id, instrumentId)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListPostgresRepository) Get(userId uint64) (*models.ComparisonList, error) {
	query := `select * from store.comparisonLists where user_id = $1;`
	comparisonListPostgres := &ComparisonListPostgres{}

	err := i.db.Get(comparisonListPostgres, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	comparisonList := &models.ComparisonList{}
	err = copier.Copy(comparisonList, comparisonListPostgres)
	if err != nil {
		return nil, err
	}

	return comparisonList, nil
}

func (i *ComparisonListPostgresRepository) GetUser(id uint64) (*models.User, error) {
	query := `select c.user_id, user_login, user_password,
					user_fio, user_date_birth, user_gender, user_is_admin
			  from store.comparisonlists c join store.users u on c.user_id = u.user_id
			  where c.user_id = $1;`

	userPostgres := &UserPostgres{}

	err := i.db.Get(userPostgres, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = copier.Copy(user, userPostgres)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *ComparisonListPostgresRepository) GetInstruments(userId uint64) ([]models.Instrument, error) {
	query := `select distinct i.instrument_id, instrument_name, instrument_price, instrument_material,
       instrument_type, instrument_brand, instrument_img
			  from (store.comparisonlists c join store.comparisonlists_instruments ci on c.comparisonlist_id = ci.comparisonlist_id) as t1
   			  join store.instruments i on t1.instrument_id = i.instrument_id
			  where t1.user_id = $1;`

	var instrumentsPostgres []InstrumentPostgres
	var instruments []models.Instrument
	err := i.db.Select(&instrumentsPostgres, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range instrumentsPostgres {
		instrument := &models.Instrument{}
		err = copier.Copy(instrument, &instrumentsPostgres[i])
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, *instrument)
	}
	return instruments, nil
}
