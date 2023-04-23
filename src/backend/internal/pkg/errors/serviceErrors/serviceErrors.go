package serviceErrors

import (
	"errors"
	"fmt"
)

var (
	Invalid           = errors.New("invalid")
	InvalidPassword   = fmt.Errorf("%w password: ", Invalid)
	InvalidDiscount   = fmt.Errorf("%w discount: ", Invalid)
	InvalidInstrument = fmt.Errorf("%w instrument: ", Invalid)

	DoesNotExists               = errors.New("does not exists: ")
	UserDoesNotExists           = fmt.Errorf("user %w", DoesNotExists)
	ComparisonListDoesNotExists = fmt.Errorf("comparisonList %w", DoesNotExists)
	InstrumentDoesNotExists     = fmt.Errorf("instrument %w", DoesNotExists)
	InstrumentsDoesNotExists    = fmt.Errorf("instruments %w", DoesNotExists)
	DiscountDoesNotExists       = fmt.Errorf("discount %w", DoesNotExists)
	DiscountsDoesNotExists      = fmt.Errorf("discounts %w", DoesNotExists)

	CanNotCreate               = errors.New("can not create")
	UserCanNotCreateInstrument = fmt.Errorf("user %w instrument, need admin: ", CanNotCreate)
	UserCanNotCreateDiscount   = fmt.Errorf("user %w discount, need admin: ", CanNotCreate)

	CanNotUpdate               = errors.New("can not update")
	UserCanNotUpdateInstrument = fmt.Errorf("user %w instrument, need admin: ", CanNotUpdate)
	UserCanNotUpdateDiscount   = fmt.Errorf("user %w discount, need admin: ", CanNotUpdate)

	CanNotDelete               = errors.New("can not delete")
	UserCanNotDeleteInstrument = fmt.Errorf("user %w instrument, need admin: ", CanNotDelete)
	UserCanNotDeleteDiscount   = fmt.Errorf("user %w discount, need admin: ", CanNotDelete)

	AlreadyExists     = errors.New("already exists: ")
	UserAlreadyExists = fmt.Errorf("user %w", AlreadyExists)

	WrongType         = errors.New("wrong type: ")
	DiscountWrongType = fmt.Errorf("discount %w", WrongType)

	CreateFailed               = errors.New("create failed: ")
	UserCreateFailed           = fmt.Errorf("user %w", CreateFailed)
	ComparisonListCreateFailed = fmt.Errorf("comparisonList %w", CreateFailed)

	GetFailed               = errors.New("get failed: ")
	UserGetFailed           = fmt.Errorf("user %w", GetFailed)
	ComparisonListGetFailed = fmt.Errorf("comparisonList %w", GetFailed)

	UpdateFailed               = errors.New("update failed: ")
	ComparisonListUpdateFailed = fmt.Errorf("comparisonList %w", UpdateFailed)
)
