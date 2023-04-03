package postgres_repository

import (
	"backend/config"
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/repository"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

type instrumentPostgresRepositoryFields struct {
	db     *sql.DB
	config config.Config
}

func createInstrumentPostgresRepository(fields *instrumentPostgresRepositoryFields) repository.InstrumentRepository {
	dbx := sqlx.NewDb(fields.db, "pgx")

	return NewInstrumentPostgresRepository(dbx)
}

func createInstrumentPostgresRepositoryFields() *instrumentPostgresRepositoryFields {
	fields := new(instrumentPostgresRepositoryFields)
	err := fields.config.ParseConfig("config.json", "../../../config")
	if err != nil {
		return nil
	}
	fields.db, err = fields.config.Postgres.InitDB()
	if err != nil {
		return nil
	}
	return fields
}

var testInstrumentPostgresRepositoryDeleteSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 0},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentPostgresRepositoryDeleteFailed = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "instrument does not exists",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 80000000},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestInstrumentPostgresRepositoryDelete(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentPostgresRepositoryDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			instrumentRepository.Create(&models.Instrument{InstrumentId: 0})
			err := instrumentRepository.Delete(tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			err := instrumentRepository.Delete(tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentPostgresRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		instrument *models.Instrument
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrument *models.Instrument
		}{instrument: &models.Instrument{InstrumentId: 8000, Name: "testing", Price: 3000}},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentPostgresRepositoryCreateFailed = []struct {
	TestName  string
	InputData struct {
		instrument *models.Instrument
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "instrument with that id already exists",
		InputData: struct {
			instrument *models.Instrument
		}{instrument: &models.Instrument{InstrumentId: 1}},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestInstrumentPostgresRepositoryCreate(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentPostgresRepositoryCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			err := instrumentRepository.Create(tt.InputData.instrument)
			instrumentRepository.Delete(tt.InputData.instrument.InstrumentId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryCreateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			err := instrumentRepository.Create(tt.InputData.instrument)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentPostgresRepositoryUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId   uint64
		fieldsToUpdate models.InstrumentFieldsToUpdate
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId   uint64
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 0, fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "qt", models.InstrumentFieldPrice: 3000}},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentPostgresRepositoryUpdateFailed = []struct {
	TestName  string
	InputData struct {
		instrumentId   uint64
		fieldsToUpdate models.InstrumentFieldsToUpdate
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "instrument does not exists",
		InputData: struct {
			instrumentId   uint64
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 218939393, fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "qt", models.InstrumentFieldPrice: rand.Intn(40404040)}},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestInstrumentPostgresRepositoryUpdate(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentPostgresRepositoryUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			instrumentRepository.Create(&models.Instrument{InstrumentId: 0})

			rand.Seed(time.Now().Unix())
			err := instrumentRepository.Update(tt.InputData.instrumentId, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createInstrumentPostgresRepositoryFields()

			instrumentRepository := createInstrumentPostgresRepository(fields)

			err := instrumentRepository.Update(tt.InputData.instrumentId, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}
