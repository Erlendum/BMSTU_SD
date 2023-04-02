// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/comparisonList.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	models "backend/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockComparisonListRepository is a mock of ComparisonListRepository interface.
type MockComparisonListRepository struct {
	ctrl     *gomock.Controller
	recorder *MockComparisonListRepositoryMockRecorder
}

// MockComparisonListRepositoryMockRecorder is the mock recorder for MockComparisonListRepository.
type MockComparisonListRepositoryMockRecorder struct {
	mock *MockComparisonListRepository
}

// NewMockComparisonListRepository creates a new mock instance.
func NewMockComparisonListRepository(ctrl *gomock.Controller) *MockComparisonListRepository {
	mock := &MockComparisonListRepository{ctrl: ctrl}
	mock.recorder = &MockComparisonListRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComparisonListRepository) EXPECT() *MockComparisonListRepositoryMockRecorder {
	return m.recorder
}

// AddInstrument mocks base method.
func (m *MockComparisonListRepository) AddInstrument(id uint64, instrumentId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddInstrument", id, instrumentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddInstrument indicates an expected call of AddInstrument.
func (mr *MockComparisonListRepositoryMockRecorder) AddInstrument(id, instrument interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddInstrument", reflect.TypeOf((*MockComparisonListRepository)(nil).AddInstrument), id, instrument)
}

// Create mocks base method.
func (m *MockComparisonListRepository) Create(comparisonList *models.ComparisonList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", comparisonList)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockComparisonListRepositoryMockRecorder) Create(comparisonList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockComparisonListRepository)(nil).Create), comparisonList)
}

// DeleteInstrument mocks base method.
func (m *MockComparisonListRepository) DeleteInstrument(id, instrumentId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInstrument", id, instrumentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteInstrument indicates an expected call of DeleteInstrument.
func (mr *MockComparisonListRepositoryMockRecorder) DeleteInstrument(id, instrumentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInstrument", reflect.TypeOf((*MockComparisonListRepository)(nil).DeleteInstrument), id, instrumentId)
}

// Get mocks base method.
func (m *MockComparisonListRepository) Get(userId uint64) (*models.ComparisonList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", userId)
	ret0, _ := ret[0].(*models.ComparisonList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockComparisonListRepositoryMockRecorder) Get(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockComparisonListRepository)(nil).Get), userId)
}

// GetInstruments mocks base method.
func (m *MockComparisonListRepository) GetInstruments(userId uint64) ([]models.Instrument, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstruments", userId)
	ret0, _ := ret[0].([]models.Instrument)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstruments indicates an expected call of GetInstruments.
func (mr *MockComparisonListRepositoryMockRecorder) GetInstruments(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstruments", reflect.TypeOf((*MockComparisonListRepository)(nil).GetInstruments), userId)
}

// GetUser mocks base method.
func (m *MockComparisonListRepository) GetUser(id uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockComparisonListRepositoryMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockComparisonListRepository)(nil).GetUser), id)
}
