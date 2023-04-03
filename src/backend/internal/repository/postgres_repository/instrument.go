package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/queries"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type instrumentPostgresRepository struct {
	db *sqlx.DB
}

func NewInstrumentPostgresRepository(db *sqlx.DB) repository.InstrumentRepository {
	return &instrumentPostgresRepository{db: db}
}

func (i *instrumentPostgresRepository) Create(instrument *models.Instrument) error {
	query := `insert into store.instruments (instrument_id, instrument_name, instrument_price, instrument_material,
											 instrument_type, instrument_brand, instrument_img) values
											 ($1, $2, $3, $4, $5, $6, $7);`
	_, err := i.db.Exec(query, instrument.InstrumentId, instrument.Name, instrument.Price, instrument.Material, instrument.Type, instrument.Brand, instrument.Img)
	if err != nil {
		return err
	}
	return nil
}

func (i *instrumentPostgresRepository) instrumentFieldToDBField(field models.InstrumentField) string {
	switch field {
	case models.InstrumentFieldName:
		return "instrument_name"
	case models.InstrumentFieldPrice:
		return "instrument_price"
	case models.InstrumentFieldMaterial:
		return "instrument_material"
	case models.InstrumentFieldType:
		return "instrument_type"
	case models.InstrumentFieldBrand:
		return "instrument_brand"
	case models.InstrumentFieldImg:
		return "instrument_img"
	}
	return ""
}

func (i *instrumentPostgresRepository) Update(id uint64, fieldsToUpdate models.InstrumentFieldsToUpdate) error {
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		updateFields[i.instrumentFieldToDBField(key)] = value
	}

	query, fields := queries.CreateUpdateQuery("store.instruments", updateFields)

	fields = append(fields, id)
	query += ` where instrument_id = $` + strconv.Itoa(len(fields)) + ";"

	res, err := i.db.Exec(query, fields...)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (i *instrumentPostgresRepository) Delete(id uint64) error {
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

func (i *instrumentPostgresRepository) Get(id uint64) (*models.Instrument, error) {
	query := `select * from store.instruments where instrument_id = $1`
	instrument := &models.Instrument{}

	err := i.db.Get(instrument, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	return instrument, nil
}

func (i *instrumentPostgresRepository) GetList() ([]models.Instrument, error) {
	query := `select * from store.instruments;`

	var instruments []models.Instrument
	err := i.db.Select(&instruments, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	return instruments, nil
}
