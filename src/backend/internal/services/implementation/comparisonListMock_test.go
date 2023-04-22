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

type comparisonListServiceFields struct {
	comparisonListRepositoryMock *mock_repository.MockComparisonListRepository
	instrumentRepositoryMock     *mock_repository.MockInstrumentRepository
}

func createComparisonListServiceFields(controller *gomock.Controller) *comparisonListServiceFields {
	fields := new(comparisonListServiceFields)

	fields.comparisonListRepositoryMock = mock_repository.NewMockComparisonListRepository(controller)
	fields.instrumentRepositoryMock = mock_repository.NewMockInstrumentRepository(controller)

	return fields
}

func createComparisonListService(fields *comparisonListServiceFields) services.ComparisonListService {
	return NewComparisonListServiceImplementation(fields.comparisonListRepositoryMock, fields.instrumentRepositoryMock, logger.New(""))
}

var testAddInstrumentSuccess = []struct {
	TestName  string
	InputData struct {
		id           uint64
		instrumentId uint64
	}
	Prepare     func(fields *comparisonListServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().AddInstrument(uint64(1), uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testAddInstrumentFailed = []struct {
	TestName  string
	InputData struct {
		id           uint64
		instrumentId uint64
	}
	Prepare     func(fields *comparisonListServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "comparisonList does not exists",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.ComparisonListDoesNotExists)
		},
	},
	{
		TestName: "comparisonList repository internal error",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "instrument does not exists",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestComparisonListServiceImplementationAddInstrument(t *testing.T) {
	t.Parallel()

	for _, tt := range testAddInstrumentSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createComparisonListServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createComparisonListService(fields)

			err := instrumentService.AddInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testAddInstrumentFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createComparisonListServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createComparisonListService(fields)

			err := instrumentService.AddInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}

var testDeleteInstrumentSuccess = []struct {
	TestName  string
	InputData struct {
		id           uint64
		instrumentId uint64
	}
	Prepare     func(fields *comparisonListServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.Instrument{InstrumentId: 1}, nil)
			fields.comparisonListRepositoryMock.EXPECT().DeleteInstrument(uint64(1), uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDeleteInstrumentFailed = []struct {
	TestName  string
	InputData struct {
		id           uint64
		instrumentId uint64
	}
	Prepare     func(fields *comparisonListServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "comparisonList does not exists",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.ComparisonListDoesNotExists)
		},
	},
	{
		TestName: "comparisonList repository internal error",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "instrument does not exists",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, serviceErrors.InstrumentDoesNotExists)
		},
	},
	{
		TestName: "instrument repository internal error",
		InputData: struct {
			id           uint64
			instrumentId uint64
		}{id: 1, instrumentId: 1},
		Prepare: func(fields *comparisonListServiceFields) {
			fields.comparisonListRepositoryMock.EXPECT().Get(uint64(1)).Return(&models.ComparisonList{ComparisonListId: 1}, nil)
			fields.instrumentRepositoryMock.EXPECT().Get(uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
}

func TestComparisonListServiceImplementationDeleteInstrument(t *testing.T) {
	t.Parallel()

	for _, tt := range testDeleteInstrumentSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createComparisonListServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createComparisonListService(fields)

			err := instrumentService.DeleteInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDeleteInstrumentFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createComparisonListServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createComparisonListService(fields)

			err := instrumentService.DeleteInstrument(tt.InputData.id, tt.InputData.instrumentId)

			tt.CheckOutput(t, err)
		})
	}
}
