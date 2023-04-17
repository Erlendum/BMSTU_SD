package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/repository/postgres_repository"
	"backend/internal/services"
	"github.com/stretchr/testify/require"
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

func createUserServiceFieldsPostgres() *userServiceFieldsPostgres {
	fields := new(userServiceFieldsPostgres)
	fieldsPostgres := postgres_repository.CreatePostgresRepositoryFields("config.json", "../../../config")

	calcDiscountServiceFields := new(calcDiscountServiceFieldsForUserPostgres)

	discountRepository := postgres_repository.CreateDiscountPostgresRepository(fieldsPostgres)
	comparisonListRepository := postgres_repository.CreateComparisonListPostgresRepository(fieldsPostgres)
	userRepository := postgres_repository.CreateUserPostgresRepository(fieldsPostgres)

	calcDiscountServiceFields.discountRepository = &discountRepository

	fields.comparisonListRepository = &comparisonListRepository
	fields.userRepository = &userRepository
	fields.calcDiscountService = NewCalcDiscountServiceImplementation(discountRepository)
	fields.discountRepository = &discountRepository
	return fields
}

func createUserServicePostgres(fields *userServiceFieldsPostgres) services.UserService {
	return NewUserServiceImplementation(*fields.userRepository, *fields.comparisonListRepository, fields.calcDiscountService)
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
			require.Equal(t, comparisonList, &models.ComparisonList{ComparisonListId: 1, UserId: 1, Amount: 1, TotalPrice: 4050})
			require.Equal(t, instruments, []models.Instrument{{
				InstrumentId: 1,
				Name:         "KALA KA-15S Kala Mahogany Soprano Ukulele No Binding",
				Price:        4050,
				Material:     "Сосна",
				Type:         "Укулеле",
				Brand:        "KALA",
				Img:          "https://www.muztorg.ru/files/sized/f250x250/vop/cd1/73u/00g/oc0/sso/kc0/408/vopcd173u00goc0ssokc0408.jpg",
			}})
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
}
