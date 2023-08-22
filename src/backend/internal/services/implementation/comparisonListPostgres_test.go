package servicesImplementation

import (
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

type comparisonListServiceFieldsPostgres struct {
	comparisonListRepository *repository.ComparisonListRepository
	instrumentRepository     *repository.InstrumentRepository
}

var comparisonListDbContainer testcontainers.Container

func createComparisonListServiceFieldsPostgres() *comparisonListServiceFieldsPostgres {
	fields := new(comparisonListServiceFieldsPostgres)

	var db *sql.DB
	comparisonListDbContainer, db = postgres_repository.SetupTestDatabase("../../repository/postgres_repository/migrations/000001_create_init_tables.up.sql")

	repositoryFields := new(postgres_repository.PostgresRepositoryFields)
	repositoryFields.Db = db
	instrumentRepository := postgres_repository.CreateInstrumentPostgresRepository(repositoryFields)
	comparisonListRepository := postgres_repository.CreateComparisonListPostgresRepository(repositoryFields)

	fields.instrumentRepository = &instrumentRepository
	fields.comparisonListRepository = &comparisonListRepository
	return fields
}

func createComparisonListServicePostgres(fields *comparisonListServiceFieldsPostgres) services.ComparisonListService {
	return NewComparisonListServiceImplementation(*fields.comparisonListRepository, *fields.instrumentRepository, logger.New("", ""))
}

var testAddInstrumentPostgresSuccess = []struct {
	TestName  string
	InputData struct {
		id           uint64
		instrumentId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestComparisonListServiceImplementationAddInstrumentPostgres(t *testing.T) {
	for _, tt := range testAddInstrumentPostgresSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createComparisonListServiceFieldsPostgres()

			comparisonListService := createComparisonListServicePostgres(fields)

			comparisonListService.DeleteInstrument(tt.InputData.id, tt.InputData.instrumentId)
			err := comparisonListService.AddInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
	err := comparisonListDbContainer.Terminate(context.Background())
	if err != nil {
		return
	}
}
