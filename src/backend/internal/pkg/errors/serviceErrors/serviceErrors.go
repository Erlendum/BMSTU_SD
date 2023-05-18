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
	OrderDoesNotExists          = fmt.Errorf("order %w", DoesNotExists)
	OrdersDoesNotExists         = fmt.Errorf("orders %w", DoesNotExists)

	CanNotCreate               = errors.New("can not create")
	UserCanNotCreateInstrument = fmt.Errorf("user %w instrument, need admin: ", CanNotCreate)
	UserCanNotCreateDiscount   = fmt.Errorf("user %w discount, need admin: ", CanNotCreate)

	CanNotUpdate               = errors.New("can not update")
	UserCanNotUpdateInstrument = fmt.Errorf("user %w instrument, need admin: ", CanNotUpdate)
	UserCanNotUpdateDiscount   = fmt.Errorf("user %w discount, need admin: ", CanNotUpdate)
	UserCanNotUpdateOrder      = fmt.Errorf("user %w order, need admin: ", CanNotUpdate)

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
	InstrumentCreateFailed     = fmt.Errorf("instrument %w", CreateFailed)
	DiscountCreateFailed       = fmt.Errorf("discount %w", CreateFailed)
	DiscountForAllCreateFailed = fmt.Errorf("discount for all %w", CreateFailed)
	OrderCreateFailed          = fmt.Errorf("order %w", CreateFailed)

	GetFailed                 = errors.New("get failed: ")
	UserGetFailed             = fmt.Errorf("user %w", GetFailed)
	ComparisonListGetFailed   = fmt.Errorf("comparisonList %w", GetFailed)
	InstrumentGetFailed       = fmt.Errorf("instrument %w", GetFailed)
	InstrumentsListGetFailed  = fmt.Errorf("instruments list %w", GetFailed)
	DiscountGetFailed         = fmt.Errorf("discount %w", GetFailed)
	DiscountsListGetFailed    = fmt.Errorf("discounts list %w", GetFailed)
	OrdersListGetFailed       = fmt.Errorf("orders list %w", GetFailed)
	OrdersListForAllGetFailed = fmt.Errorf("orders list for all %w", GetFailed)

	UpdateFailed               = errors.New("update failed: ")
	ComparisonListUpdateFailed = fmt.Errorf("comparisonList %w", UpdateFailed)
	InstrumentUpdateFailed     = fmt.Errorf("instrument %w", UpdateFailed)
	DiscountUpdateFailed       = fmt.Errorf("discount %w", UpdateFailed)
	OrderUpdateFailed          = fmt.Errorf("order %w", UpdateFailed)

	DeleteFailed           = errors.New("delete failed: ")
	InstrumentDeleteFailed = fmt.Errorf("instrument %w", DeleteFailed)
	DiscountDeleteFailed   = fmt.Errorf("delete %w", DeleteFailed)

	ComparisonListAddInstrumentFailed    = errors.New("comparisonList add instrument failed: ")
	ComparisonListDeleteInstrumentFailed = errors.New("comparisonList delete instrument failed: ")
)
