package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/hasher"
	"backend/internal/pkg/hasher/implementation"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type calcDiscountServiceFieldsForUser struct {
	discountRepositoryMock *mock_repository.MockDiscountRepository
}

type userServiceFields struct {
	comparisonListRepositoryMock *mock_repository.MockComparisonListRepository
	userRepositoryMock           *mock_repository.MockUserRepository
	discountRepositoryMock       *mock_repository.MockDiscountRepository
	calcDiscountService          services.CalcDiscountService
	hasher                       hasher.Hasher
}

func createUserServiceFields(controller *gomock.Controller) *userServiceFields {
	calcDiscountServiceFields := new(calcDiscountServiceFieldsForUser)
	calcDiscountServiceFields.discountRepositoryMock = mock_repository.NewMockDiscountRepository(controller)
	fields := new(userServiceFields)

	fields.comparisonListRepositoryMock = mock_repository.NewMockComparisonListRepository(controller)
	fields.userRepositoryMock = mock_repository.NewMockUserRepository(controller)
	fields.calcDiscountService = NewCalcDiscountServiceImplementation(calcDiscountServiceFields.discountRepositoryMock)
	fields.discountRepositoryMock = calcDiscountServiceFields.discountRepositoryMock
	fields.hasher = &implementation.BcryptHasher{}
	return fields
}

func createUserService(fields *userServiceFields) services.UserService {
	return NewUserServiceImplementation(fields.userRepositoryMock, fields.comparisonListRepositoryMock, fields.calcDiscountService)
}

var testUserCreateFailed = []struct {
	TestName  string
	InputData struct {
		user     *models.User
		password string
	}
	Prepare     func(fields *userServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user already exists",
		InputData: struct {
			user     *models.User
			password string
		}{user: &models.User{UserId: 1, Login: "login1"}, password: "123"},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{UserId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserAlreadyExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			user     *models.User
			password string
		}{user: &models.User{UserId: 1, Login: "login1"}, password: "123"},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestUserServiceImplementationCreate(t *testing.T) {
	t.Parallel()

	for _, tt := range testUserCreateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createUserServiceFields(ctrl)
			tt.Prepare(fields)

			userService := createUserService(fields)

			err := userService.Create(tt.InputData.user, tt.InputData.password)

			tt.CheckOutput(t, err)
		})
	}
}

var testUserGetSuccess = []struct {
	TestName  string
	InputData struct {
		login    string
		password string
	}
	Prepare     func(fields *userServiceFields)
	CheckOutput func(t *testing.T, user *models.User, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			login    string
			password string
		}{login: "login1", password: "123"},
		Prepare: func(fields *userServiceFields) {
			hashPassword, _ := fields.hasher.GetHash("123")
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{Login: "login1", Password: string(hashPassword)}, nil)
		},
		CheckOutput: func(t *testing.T, user *models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, user.Login, (&models.User{Login: "login1"}).Login)
		},
	},
}

var testUserGetFailed = []struct {
	TestName  string
	InputData struct {
		login    string
		password string
	}
	Prepare     func(fields *userServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user doest not exists",
		InputData: struct {
			login    string
			password string
		}{login: "login1", password: "123"},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "invalid password",
		InputData: struct {
			login    string
			password string
		}{login: "login1", password: "1234"},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{Login: "login1", Password: "123"}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InvalidPassword)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			login    string
			password string
		}{login: "login1", password: "1234"},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestUserServiceImplementationGet(t *testing.T) {
	t.Parallel()

	for _, tt := range testUserGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createUserServiceFields(ctrl)
			tt.Prepare(fields)

			userService := createUserService(fields)

			user, err := userService.Get(tt.InputData.login, tt.InputData.password)

			tt.CheckOutput(t, user, err)
		})
	}

	for _, tt := range testUserGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createUserServiceFields(ctrl)
			tt.Prepare(fields)

			userService := createUserService(fields)

			_, err := userService.Get(tt.InputData.login, tt.InputData.password)

			tt.CheckOutput(t, err)
		})
	}
}

var testUserGetComparisonListSuccess = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *userServiceFields)
	CheckOutput func(t *testing.T, comparisonList *models.ComparisonList, instruments []models.Instrument, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().GetInstruments(uint64(1)).Return([]models.Instrument{{InstrumentId: 1}}, nil)
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, comparisonList *models.ComparisonList, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, comparisonList, &models.ComparisonList{ComparisonListId: 1})
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1}})
		},
	},

	{
		TestName: "no instruments",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().GetInstruments(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, comparisonList *models.ComparisonList, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, comparisonList, &models.ComparisonList{ComparisonListId: 1})
			require.Equal(t, instruments, []models.Instrument(nil))
		},
	},
	{
		TestName: "no discounts",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().GetInstruments(uint64(1)).Return([]models.Instrument{{InstrumentId: 1}}, nil)
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, comparisonList *models.ComparisonList, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, comparisonList, &models.ComparisonList{ComparisonListId: 1})
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1}})
		},
	},
}

var testUserGetComparisonListFailed = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *userServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user doest not exists",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "comparisonList does not exists",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.ComparisonListDoesNotExists)
		},
	},
	{
		TestName: "comparisonList repository internal error 1",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "comparisonList repository internal error 2",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().GetInstruments(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "discount repository internal error",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *userServiceFields) {
			fields.userRepositoryMock.EXPECT().GetById(uint64(1)).Return(&models.User{UserId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().GetInstruments(uint64(1)).Return([]models.Instrument{{InstrumentId: 1}}, nil)
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestUserServiceImplementationGetComparisonList(t *testing.T) {
	t.Parallel()

	for _, tt := range testUserGetComparisonListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createUserServiceFields(ctrl)
			tt.Prepare(fields)

			userService := createUserService(fields)

			comparisonList, instruments, err := userService.GetComparisonList(tt.InputData.id)

			tt.CheckOutput(t, comparisonList, instruments, err)
		})
	}

	for _, tt := range testUserGetComparisonListFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createUserServiceFields(ctrl)
			tt.Prepare(fields)

			userService := createUserService(fields)

			_, _, err := userService.GetComparisonList(tt.InputData.id)

			tt.CheckOutput(t, err)
		})
	}
}
