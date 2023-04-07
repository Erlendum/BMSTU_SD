package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var testDiscountPostgresRepositoryDeleteSuccess = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId uint64
		}{discountId: 0},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountPostgresRepositoryDeleteFailed = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId uint64
		}{discountId: 80000000},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestDiscountPostgresRepositoryDelete(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			discountRepository.Create(&models.Discount{DiscountId: 0, InstrumentId: 1, UserId: 1})
			err := discountRepository.Delete(tt.InputData.discountId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountPostgresRepositoryDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			err := discountRepository.Delete(tt.InputData.discountId)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountPostgresRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		discount *models.Discount
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discount *models.Discount
		}{discount: &models.Discount{DiscountId: 1, InstrumentId: 1, UserId: 1}},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountPostgresRepositoryCreateFailed = []struct {
	TestName  string
	InputData struct {
		discount *models.Discount
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount with that id already exists",
		InputData: struct {
			discount *models.Discount
		}{discount: &models.Discount{DiscountId: 1, InstrumentId: 1, UserId: 1}},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestDiscountPostgresRepositoryCreate(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)
			discountRepository.Delete(tt.InputData.discount.InstrumentId)
			err := discountRepository.Create(tt.InputData.discount)
			discountRepository.Delete(tt.InputData.discount.InstrumentId)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountPostgresRepositoryCreateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			discountRepository.Create(tt.InputData.discount)
			err := discountRepository.Create(tt.InputData.discount)
			discountRepository.Delete(tt.InputData.discount.InstrumentId)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountPostgresRepositoryUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		discountId     uint64
		fieldsToUpdate models.DiscountFieldsToUpdate
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId     uint64
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 0, fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldAmount: 3000, models.DiscountFieldType: "Процентная"}},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDiscountPostgresRepositoryUpdateFailed = []struct {
	TestName  string
	InputData struct {
		discountId     uint64
		fieldsToUpdate models.DiscountFieldsToUpdate
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId     uint64
			fieldsToUpdate models.DiscountFieldsToUpdate
		}{discountId: 2122, fieldsToUpdate: map[models.DiscountField]any{models.DiscountFieldAmount: 3000, models.DiscountFieldType: "Процентная"}},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestDiscountPostgresRepositoryUpdate(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			discountRepository.Create(&models.Discount{DiscountId: 0, UserId: 1, InstrumentId: 1})

			rand.Seed(time.Now().Unix())
			err := discountRepository.Update(tt.InputData.discountId, tt.InputData.fieldsToUpdate)

			discountRepository.Delete(tt.InputData.discountId)
			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testDiscountPostgresRepositoryUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			err := discountRepository.Update(tt.InputData.discountId, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountPostgresRepositoryGetSuccess = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	CheckOutput func(t *testing.T, discount *models.Discount, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			discountId uint64
		}{discountId: 0},
		CheckOutput: func(t *testing.T, discount *models.Discount, err error) {
			require.NoError(t, err)
			require.Equal(t, discount, &models.Discount{DiscountId: 0, InstrumentId: 1, UserId: 1})
		},
	},
}

var testDiscountPostgresRepositoryGetFailed = []struct {
	TestName  string
	InputData struct {
		discountId uint64
	}
	CheckOutput func(t *testing.T, err error)
}{

	{
		TestName: "discount does not exists",
		InputData: struct {
			discountId uint64
		}{discountId: 218939393},
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestDiscountPostgresRepositoryGet(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			discountRepository.Create(&models.Discount{DiscountId: 0, InstrumentId: 1, UserId: 1})

			discount, err := discountRepository.Get(tt.InputData.discountId)

			discountRepository.Delete(tt.InputData.discountId)
			tt.CheckOutput(t, discount, err)
		})
	}

	for _, tt := range testDiscountPostgresRepositoryGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			instrumentRepository := CreateInstrumentPostgresRepository(fields)

			_, err := instrumentRepository.Get(tt.InputData.discountId)

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountPostgresRepositoryGetListSuccess = []struct {
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

func TestDiscountPostgresRepositoryGetList(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryGetListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			_, err := discountRepository.GetList()

			tt.CheckOutput(t, err)
		})
	}
}

var testDiscountPostgresRepositoryGetSpecificListSuccess = []struct {
	TestName  string
	InputData struct {
		instrumentId uint64
		userId       uint64
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			instrumentId uint64
			userId       uint64
		}{instrumentId: 1, userId: 1},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestDiscountPostgresRepositoryGetSpecificList(t *testing.T) {
	for _, tt := range testDiscountPostgresRepositoryGetSpecificListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := CreatePostgresRepositoryFields()

			discountRepository := CreateDiscountPostgresRepository(fields)

			_, err := discountRepository.GetSpecificList(tt.InputData.instrumentId, tt.InputData.userId)

			tt.CheckOutput(t, err)
		})
	}
}
