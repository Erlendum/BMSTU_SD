package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/errors/serviceErrors"
	"backend/internal/pkg/logger"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type instrumentServiceFields struct {
	instrumentRepositoryMock *mock_repository.MockInstrumentRepository
	userRepositoryMock       *mock_repository.MockUserRepository
}

func createInstrumentServiceFields(controller *gomock.Controller) *instrumentServiceFields {
	fields := new(instrumentServiceFields)

	fields.instrumentRepositoryMock = mock_repository.NewMockInstrumentRepository(controller)
	fields.userRepositoryMock = mock_repository.NewMockUserRepository(controller)

	return fields
}

func createInstrumentService(fields *instrumentServiceFields) services.InstrumentService {
	return NewInstrumentServiceImplementation(fields.instrumentRepositoryMock, fields.userRepositoryMock, logger.New("", ""))
}

var testInstrumentCreateSuccess = []struct {
	TestName  string
	InputData struct {
		instrument *models.Instrument
		login      string
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrument *models.Instrument
			login      string
		}{instrument: &models.Instrument{InstrumentId: 1}, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.instrumentRepositoryMock.EXPECT().Create(&models.Instrument{InstrumentId: 1}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentCreateFailed = []struct {
	TestName  string
	InputData struct {
		instrument *models.Instrument
		login      string
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			instrument *models.Instrument
			login      string
		}{instrument: &models.Instrument{InstrumentId: 1}, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			instrument *models.Instrument
			login      string
		}{instrument: &models.Instrument{InstrumentId: 1}, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not create instrument",
		InputData: struct {
			instrument *models.Instrument
			login      string
		}{instrument: &models.Instrument{InstrumentId: 1}, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotCreateInstrument)
		},
	},
}

func TestInstrumentServiceImplementationCreate(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Create(tt.InputData.instrument, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentCreateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Create(tt.InputData.instrument, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId   uint64
		login          string
		fieldsToUpdate models.InstrumentFieldsToUpdate
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Update(uint64(1), map[models.InstrumentField]any{models.InstrumentFieldName: "p"}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentUpdateFailed = []struct {
	TestName  string
	InputData struct {
		instrumentId   uint64
		login          string
		fieldsToUpdate models.InstrumentFieldsToUpdate
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "instrument does not exists",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not update instrument",
		InputData: struct {
			instrumentId   uint64
			login          string
			fieldsToUpdate models.InstrumentFieldsToUpdate
		}{instrumentId: 1, login: "login1", fieldsToUpdate: map[models.InstrumentField]any{models.InstrumentFieldName: "p"}},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotUpdateInstrument)
		},
	},
}

func TestInstrumentServiceImplementationUpdate(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Update(tt.InputData.instrumentId, tt.InputData.login, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Update(tt.InputData.instrumentId, tt.InputData.login, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentDeleteSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
		login        string
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: true}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Delete(uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testInstrumentDeleteFailed = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
		login        string
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "user does not exists",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.ObjectDoesNotExists)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserDoesNotExists)
		},
	},
	{
		TestName: "user repository internal error",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.userRepositoryMock.EXPECT().Get("login1").Return(nil, repositoryErrors.InternalRepositoryError)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "instrument does not exists",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "user can not delete instrument",
		InputData: struct {
			instrumentId uint64
			login        string
		}{instrumentId: 1, login: "login1"},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.userRepositoryMock.EXPECT().Get("login1").Return(&models.User{IsAdmin: false}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.UserCanNotDeleteInstrument)
		},
	},
}

func TestInstrumentServiceImplementationDelete(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Delete(tt.InputData.instrumentId, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testInstrumentDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			err := instrumentService.Delete(tt.InputData.instrumentId, tt.InputData.login)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentGetSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, instrument *models.Instrument, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 1},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
		},
		CheckOutput: func(t *testing.T, instrument *models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instrument, &models.Instrument{InstrumentId: 1})
		},
	},
}

var testInstrumentGetFailed = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "instrument does not exists",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 1},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
			instrumentId uint64
		}{instrumentId: 1},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestInstrumentServiceImplementationGet(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			instrument, err := instrumentService.Get(tt.InputData.instrumentId)

			tt.CheckOutput(t, instrument, err)
		})
	}

	for _, tt := range testInstrumentGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			_, err := instrumentService.Get(tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}

var testInstrumentGetListSuccess = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, instruments []models.Instrument, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
		}{},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().GetList().Return([]models.Instrument{{InstrumentId: 1}, {InstrumentId: 2}}, nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1}, {InstrumentId: 2}})
		},
	},
}

var testInstrumentGetListFailed = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *instrumentServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "instrument does not exists",
		InputData: struct {
		}{},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().GetList().Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentsDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
		}{},
		Prepare: func(fields *instrumentServiceFields) {
			fields.instrumentRepositoryMock.EXPECT().GetList().Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestInstrumentServiceImplementationGetList(t *testing.T) {
	t.Parallel()

	for _, tt := range testInstrumentGetListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			instruments, err := instrumentService.GetList()

			tt.CheckOutput(t, instruments, err)
		})
	}

	for _, tt := range testInstrumentGetListFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createInstrumentServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createInstrumentService(fields)

			_, err := instrumentService.GetList()

			tt.CheckOutput(t, err)
		})
	}
}
