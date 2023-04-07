package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type calcDiscountServiceFields struct {
	discountRepositoryMock *mock_repository.MockDiscountRepository
}

func createCalcDiscountServiceFields(controller *gomock.Controller) *calcDiscountServiceFields {
	fields := new(calcDiscountServiceFields)

	fields.discountRepositoryMock = mock_repository.NewMockDiscountRepository(controller)

	return fields
}

func createCalcDiscountService(fields *calcDiscountServiceFields) services.CalcDiscountService {
	return NewCalcDiscountServiceImplementation(fields.discountRepositoryMock)
}

var testCalcDiscountSuccess = []struct {
	TestName  string
	InputData struct {
		user        *models.User
		instruments []models.Instrument
	}
	Prepare     func(fields *calcDiscountServiceFields)
	CheckOutput func(t *testing.T, instruments []models.Instrument, err error)
}{
	{
		TestName: "one discount of type percent",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{{InstrumentId: 1, UserId: 1, Type: "Процентная", Amount: 10, DateBegin: time.Now(), DateEnd: time.Now().AddDate(0, 1, 1)}}, nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 90}})
		},
	},
	{
		TestName: "several discounts of types percent and birth",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Процентная", Amount: 20, DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
					{InstrumentId: 1, UserId: 1, Type: "Именинная", Amount: 30, DateBegin: time.Now().AddDate(0, 0, -1), DateEnd: time.Now().AddDate(0, 0, 1)}}, nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 70}})
		},
	},
	{
		TestName: "one discount of type celebration",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Мужской", Amount: 45, DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)}}, nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 55}})
		},
	},
	{
		TestName: "several discounts of types celebration, percent and birth",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Мужской", Amount: 21, DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
					{InstrumentId: 1, UserId: 1, Type: "Процентная", Amount: 20, DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
					{InstrumentId: 1, UserId: 1, Type: "Именинная", Amount: 30, DateBegin: time.Now().AddDate(0, 0, -1), DateEnd: time.Now().AddDate(0, 0, 1)}},
				nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 70}})
		},
	},
	{
		TestName: "one discount of type gift",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Подарочная 1 3", DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
				},
				nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 100},
				{InstrumentId: 1, Price: 0}, {InstrumentId: 1, Price: 0}, {InstrumentId: 1, Price: 0}})
		},
	},
	{
		TestName: "invalid number of params discount of type gift",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Подарочная p 3 1 2", DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
				},
				nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 100}})
		},
	},
	{
		TestName: "invalid param n discount of type gift",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Подарочная p 3", DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
				},
				nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 100}})
		},
	},
	{
		TestName: "invalid param m discount of type gift",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1, Gender: "Мужской"}, instruments: []models.Instrument{{InstrumentId: 1, Price: 100}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(
				[]models.Discount{
					{InstrumentId: 1, UserId: 1, Type: "Подарочная 1 m", DateBegin: time.Now().AddDate(0, -1, 0), DateEnd: time.Now().AddDate(0, 0, 1)},
				},
				nil)
		},
		CheckOutput: func(t *testing.T, instruments []models.Instrument, err error) {
			require.NoError(t, err)
			require.Equal(t, instruments, []models.Instrument{{InstrumentId: 1, Price: 100}})
		},
	},
}

var testCalcDiscountFailed = []struct {
	TestName  string
	InputData struct {
		user        *models.User
		instruments []models.Instrument
	}
	Prepare     func(fields *calcDiscountServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount repository internal error",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1}, instruments: []models.Instrument{{InstrumentId: 1}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(nil, repositoryErrors.InternalRepositoryError)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.InternalRepositoryError)
		},
	},
	{
		TestName: "no discounts",
		InputData: struct {
			user        *models.User
			instruments []models.Instrument
		}{user: &models.User{UserId: 1}, instruments: []models.Instrument{{InstrumentId: 1}}},
		Prepare: func(fields *calcDiscountServiceFields) {
			fields.discountRepositoryMock.EXPECT().GetSpecificList(uint64(1), uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestCalcDiscountServiceImplementation(t *testing.T) {
	t.Parallel()

	for _, tt := range testCalcDiscountSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createCalcDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createCalcDiscountService(fields)

			instruments, err := instrumentService.CalcDiscount(tt.InputData.user, tt.InputData.instruments)

			tt.CheckOutput(t, instruments, err)
		})
	}

	for _, tt := range testCalcDiscountFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createCalcDiscountServiceFields(ctrl)
			tt.Prepare(fields)

			instrumentService := createCalcDiscountService(fields)

			_, err := instrumentService.CalcDiscount(tt.InputData.user, tt.InputData.instruments)

			tt.CheckOutput(t, err)
		})
	}
}
