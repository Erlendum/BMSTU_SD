package servicesImplementation

import (
	"backend/internal/models"
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

type calcDiscountServiceFieldsForUserPostgres struct {
	discountRepository *repository.DiscountRepository
}

type userServiceFieldsPostgres struct {
	comparisonListRepository *repository.ComparisonListRepository
	userRepository           *repository.UserRepository
	discountRepository       *repository.DiscountRepository
	calcDiscountService      services.CalcDiscountService
}

var userDbContainer testcontainers.Container

func createUserServiceFieldsPostgres() *userServiceFieldsPostgres {

	fields := new(userServiceFieldsPostgres)
	var db *sql.DB
	userDbContainer, db = postgres_repository.SetupTestDatabase("../../repository/postgres_repository/migrations/000001_create_init_tables.up.sql")

	repositoryFields := new(postgres_repository.PostgresRepositoryFields)
	repositoryFields.Db = db
	calcDiscountServiceFields := new(calcDiscountServiceFieldsForUserPostgres)

	discountRepository := postgres_repository.CreateDiscountPostgresRepository(repositoryFields)
	comparisonListRepository := postgres_repository.CreateComparisonListPostgresRepository(repositoryFields)
	userRepository := postgres_repository.CreateUserPostgresRepository(repositoryFields)

	calcDiscountServiceFields.discountRepository = &discountRepository

	fields.comparisonListRepository = &comparisonListRepository
	fields.userRepository = &userRepository
	fields.calcDiscountService = NewCalcDiscountServiceImplementation(discountRepository, logger.New("", ""))
	fields.discountRepository = &discountRepository
	return fields
}

func createUserServicePostgres(fields *userServiceFieldsPostgres) services.UserService {
	return NewUserServiceImplementation(*fields.userRepository, *fields.comparisonListRepository, fields.calcDiscountService, logger.New("", ""))
}

var testGetComparisonListPostgresSuccess = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	CheckOutput func(t *testing.T, err error, comparisonList *models.ComparisonList, instruments []models.Instrument)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id uint64
		}{id: 1},
		CheckOutput: func(t *testing.T, err error, comparisonList *models.ComparisonList, instruments []models.Instrument) {
			require.NoError(t, err)
			require.Equal(t, comparisonList, &models.ComparisonList{ComparisonListId: 1, UserId: 1, Amount: 0, TotalPrice: 0})
			require.Equal(t, instruments, []models.Instrument(nil))
		},
	},
}

func TestComparisonListServiceImplementationGetComparisonListPostgres(t *testing.T) {
	for _, tt := range testGetComparisonListPostgresSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := createUserServiceFieldsPostgres()

			userService := createUserServicePostgres(fields)

			comparisonList, instruments, err := userService.GetComparisonList(tt.InputData.id)

			tt.CheckOutput(t, err, comparisonList, instruments)
		})
	}

	err := userDbContainer.Terminate(context.Background())
	if err != nil {
		return
	}
}
