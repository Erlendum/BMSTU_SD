package servicesImplementation

import (
	"backend/internal/pkg/logger"
	"backend/internal/repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	"github.com/stretchr/testify/require"
	"testing"
)

type comparisonListServiceFieldsPostgres struct {
	comparisonListRepository *repository.ComparisonListRepository
	instrumentRepository     *repository.InstrumentRepository
}

func createComparisonListServiceFieldsPostgres() *comparisonListServiceFieldsPostgres {
	fields := new(comparisonListServiceFieldsPostgres)

	repositoryFields := postgres_repository.CreatePostgresRepositoryFields("config.json", "../../../config")
	instrumentRepository := postgres_repository.CreateInstrumentPostgresRepository(repositoryFields)
	comparisonListRepository := postgres_repository.CreateComparisonListPostgresRepository(repositoryFields)

	fields.instrumentRepository = &instrumentRepository
	fields.comparisonListRepository = &comparisonListRepository
	return fields
}

func createComparisonListServicePostgres(fields *comparisonListServiceFieldsPostgres) services.ComparisonListService {
	return NewComparisonListServiceImplementation(*fields.comparisonListRepository, *fields.instrumentRepository, logger.New(""))
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
	for _, tt := range testAddInstrumentSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createComparisonListServiceFieldsPostgres()

			comparisonListService := createComparisonListServicePostgres(fields)

			comparisonListService.DeleteInstrument(tt.InputData.id, tt.InputData.instrumentId)
			err := comparisonListService.AddInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}
