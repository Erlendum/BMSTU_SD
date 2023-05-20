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
)

type InstrumentPostgres struct {
	InstrumentId uint64 `db:"instrument_id"`
	Name         string `db:"instrument_name"`
	Price        uint64 `db:"instrument_price"`
	Material     string `db:"instrument_material"`
	Type         string `db:"instrument_type"`
	Brand        string `db:"instrument_brand"`
	Img          string `db:"instrument_img"`
}
type InstrumentPostgresRepository struct {
	db *sqlx.DB
}

func NewInstrumentPostgresRepository(db *sqlx.DB) repository.InstrumentRepository {
	return &InstrumentPostgresRepository{db: db}
}

func (i *InstrumentPostgresRepository) Create(instrument *models.Instrument) error {
	query := `insert into store.instruments (instrument_name, instrument_price, instrument_material,
											 instrument_type, instrument_brand, instrument_img) values
											 ($1, $2, $3, $4, $5, $6);`
	_, err := i.db.Exec(query, instrument.Name, instrument.Price, instrument.Material, instrument.Type, instrument.Brand, instrument.Img)
	if err != nil {
		return err
	}
	return nil
}

func (i *InstrumentPostgresRepository) instrumentFieldToDBField(field models.InstrumentField) (string, error) {
	switch field {
	case models.InstrumentFieldName:
		return "instrument_name", nil
	case models.InstrumentFieldPrice:
		return "instrument_price", nil
	case models.InstrumentFieldMaterial:
		return "instrument_material", nil
	case models.InstrumentFieldType:
		return "instrument_type", nil
	case models.InstrumentFieldBrand:
		return "instrument_brand", nil
	case models.InstrumentFieldImg:
		return "instrument_img", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *InstrumentPostgresRepository) Update(id uint64, fieldsToUpdate models.InstrumentFieldsToUpdate) error {
	if len(fieldsToUpdate) == 0 {
		return nil
	}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.instrumentFieldToDBField(key)
		if err != nil {
			return err
		}
		updateFields[field] = value
	}

	query, fields := queries.CreatePostgresUpdateQuery("store.instruments", updateFields)

	fields = append(fields, id)
	query += ` where instrument_id = $` + strconv.Itoa(len(fields)) + ";"

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

func (i *InstrumentPostgresRepository) Delete(id uint64) error {
	query := `delete from store.instruments where instrument_id = $1`
	res, err := i.db.Exec(query, id)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (i *InstrumentPostgresRepository) Get(id uint64) (*models.Instrument, error) {
	query := `select * from store.instruments where instrument_id = $1`
	instrumentPostgres := &InstrumentPostgres{}

	err := i.db.Get(instrumentPostgres, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	instrument := &models.Instrument{}
	err = copier.Copy(instrument, instrumentPostgres)
	if err != nil {
		return nil, err
	}

	return instrument, nil
}

func (i *InstrumentPostgresRepository) GetList() ([]models.Instrument, error) {
	query := `select * from store.instruments order by instrument_id;`

	var instrumentsPostgres []InstrumentPostgres
	var instruments []models.Instrument
	err := i.db.Select(&instrumentsPostgres, query)
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
