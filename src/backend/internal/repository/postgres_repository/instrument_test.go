package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

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
		}{instrumentId: 80000},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestInstrumentPostgresRepositoryDelete(t *testing.T) {
	for _, tt := range testInstrumentPostgresRepositoryDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)
			var id uint64 = 0
			fields.Db.Exec("insert into store.instruments (instrument_id, instrument_name, instrument_price, instrument_material, instrument_type, instrument_brand, instrument_img) values ($1, $2, $3, $4, $5, $6, $7)",
				id, "", 0, "", "", "", "")

			err := instrumentRepository.Delete(id)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

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

func TestInstrumentPostgresRepositoryCreate(t *testing.T) {
	for _, tt := range testInstrumentPostgresRepositoryCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			err := instrumentRepository.Create(tt.InputData.instrument)
			instrumentRepository.Delete(tt.InputData.instrument.InstrumentId)

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
	for _, tt := range testInstrumentPostgresRepositoryUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			fields.Db.Exec("insert into store.instruments (instrument_id, instrument_name, instrument_price, instrument_material, instrument_type, instrument_brand, instrument_img) values ($1, $2, $3, $4, $5, $6, $7)",
				tt.InputData.instrumentId, "", 0, "", "", "", "")

			rand.Seed(time.Now().Unix())
			err := instrumentRepository.Update(tt.InputData.instrumentId, tt.InputData.fieldsToUpdate)

			instrumentRepository.Delete(tt.InputData.instrumentId)
			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			err := instrumentRepository.Update(tt.InputData.instrumentId, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentPostgresRepositoryGetSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
	}
	CheckOutput func(t *testing.T, instrument *models.Instrument, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 0},
		CheckOutput: func(t *testing.T, instrument *models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instrument, &models.Instrument{InstrumentId: 0})
		},
	},
}

var testInstrumentPostgresRepositoryGetFailed = []struct {
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
		}{instrumentId: 218939393},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestInstrumentPostgresRepositoryGet(t *testing.T) {
	for _, tt := range testInstrumentPostgresRepositoryGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			fields.Db.Exec("insert into store.instruments (instrument_id, instrument_name, instrument_price, instrument_material, instrument_type, instrument_brand, instrument_img) values ($1, $2, $3, $4, $5, $6, $7)",
				tt.InputData.instrumentId, "", 0, "", "", "", "")

			rand.Seed(time.Now().Unix())
			instrument, err := instrumentRepository.Get(tt.InputData.instrumentId)

			instrumentRepository.Delete(0)
			tt.CheckOutput(t, instrument, err)
		})
	}

	for _, tt := range testInstrumentPostgresRepositoryGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			_, err := instrumentRepository.Get(tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentPostgresRepositoryGetListSuccess = []struct {
	TestName  string
	InputData struct {
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
		}{},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestInstrumentPostgresRepositoryGetList(t *testing.T) {
	for _, tt := range testInstrumentPostgresRepositoryGetListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields("config.json", "../../../config")

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			_, err := instrumentRepository.GetList()

			tt.CheckOutput(t, err)
		})
	}
}
