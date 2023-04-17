package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type discountServiceFields struct {
	discountRepositoryMock *mock_repository.MockDiscountRepository
	userRepositoryMock     *mock_repository.MockUserRepository
}

func createDiscountServiceFields(controller *gomock.Controller) *discountServiceFields {
	fields := new(discountServiceFields)

	fields.discountRepositoryMock = mock_repository.NewMockDiscountRepository(controller)
	fields.userRepositoryMock = mock_repository.NewMockUserRepository(controller)

	return fields
}

func createDiscountService(fields *discountServiceFields) services.DiscountService {
	return NewDiscountServiceImplementation(fields.discountRepositoryMock, fields.userRepositoryMock)
}

var testDiscountCreateSuccess = []struct {
	TestName  string
	InputData struct {
		discount *models.Discount
		login    string
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discount *models.Discount
			login    string
		}{discount: &models.Discount{DiscountId: 1, UserId: 1}, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.discountRepositoryMock.EXPECT().Create(&models.Discount{DiscountId: 1, UserId: 1}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountCreateFailed = []struct {
	TestName  string
	InputData struct {
		discount *models.Discount
		login    string
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			discount *models.Discount
			login    string
		}{discount: &models.Discount{DiscountId: 1}, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			discount *models.Discount
			login    string
		}{discount: &models.Discount{DiscountId: 1}, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not create discount",
		InputData: struct {
			discount *models.Discount
			login    string
		}{discount: &models.Discount{DiscountId: 1}, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotCreateDiscount)
		},
	},
}

func TestDiscountServiceImplementationCreate(t *testing.T) {
	t.Parallel()

	for _, tt := range testDiscountCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Create(tt.InputData.discount, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountCreateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Create(tt.InputData.discount, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		discountId     uint64
		login          string
		fieldsToUpdate models.DiscountFieldsToUpdate
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
			fields.discountRepositoryMock.EXPECT().Update(uint64(1), map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountUpdateFailed = []struct {
	TestName  string
	InputData struct {
		discountId     uint64
		login          string
		fieldsToUpdate models.DiscountFieldsToUpdate
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.DiscountDoesNotExists)
		},
	},
	{
		TestName: "discount repository internal error",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not update discount",
		InputData: struct {
			discountId     uint64
			login          string
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 1, login: "login1", fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldType: "Процентная"}},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotUpdateDiscount)
		},
	},
}

func TestDiscountServiceImplementationUpdate(t *testing.T) {
	t.Parallel()

	for _, tt := range testDiscountUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Update(tt.InputData.discountId, tt.InputData.login, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Update(tt.InputData.discountId, tt.InputData.login, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountDeleteSuccess = []struct {
	TestName  string
	InputData struct {
		discountId uint64
		login      string
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
			fields.discountRepositoryMock.EXPECT().Delete(uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountDeleteFailed = []struct {
	TestName  string
	InputData struct {
		discountId uint64
		login      string
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.DiscountDoesNotExists)
		},
	},
	{
		TestName: "discount repository internal error",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not delete discount",
		InputData: struct {
			discountId uint64
			login      string
		}{discountId: 1, login: "login1"},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotDeleteDiscount)
		},
	},
}

func TestDiscountServiceImplementationDelete(t *testing.T) {
	t.Parallel()

	for _, tt := range testDiscountDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Delete(tt.InputData.discountId, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			err := discountService.Delete(tt.InputData.discountId, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountGetSuccess = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, discount *models.Discount, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId uint64
		}{discountId: 1},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Discount{DiscountId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, discount *models.Discount, err error) {
			require.NoError(t, err)
			require.Equal(t, discount, &models.Discount{DiscountId: 1})
		},
	},
}

var testDiscountGetFailed = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId uint64
		}{discountId: 1},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.DiscountDoesNotExists)
		},
	},
	{
		TestName: "discount repository internal error",
		InputData: struct {
			discountId uint64
		}{discountId: 1},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestDiscountServiceImplementationGet(t *testing.T) {
	t.Parallel()

	for _, tt := range testDiscountGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			discount, err := discountService.Get(tt.InputData.discountId)

			tt.CheckOutput(t, discount, err)
		})
	}

	for _, tt := range testDiscountGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			_, err := discountService.Get(tt.InputData.discountId)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountGetListSuccess = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, discounts []models.Discount, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
		}{},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetList().Return([]models.Discount{{DiscountId: 1}, {DiscountId: 2}}, nil)
		},
		CheckOutput: func(t *testing.T, discounts []models.Discount, err error) {
			require.NoError(t, err)
			require.Equal(t, discounts, []models.Discount{{DiscountId: 1}, {DiscountId: 2}})
		},
	},
}

var testDiscountGetListFailed = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *discountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount does not exists",
		InputData: struct {
		}{},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetList().Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.DiscountsDoesNotExists)
		},
	},
	{
		TestName: "discount repository internal error",
		InputData: struct {
		}{},
		Prepare: func(fields *discountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetList().Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestDiscountServiceImplementationGetList(t *testing.T) {
	t.Parallel()

	for _, tt := range testDiscountGetListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			discounts, err := discountService.GetList()

			tt.CheckOutput(t, discounts, err)
		})
	}

	for _, tt := range testDiscountGetListFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			discountService := createDiscountService(fields)

			_, err := discountService.GetList()

			tt.CheckOutput(t, err)
		})
	}
}
