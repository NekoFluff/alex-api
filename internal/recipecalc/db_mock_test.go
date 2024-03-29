// Code generated by MockGen. DO NOT EDIT.
// Source: db.go

// Package recipecalc is a generated GoMock package.
package recipecalc

import (
	data "alex-api/internal/data"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// GetBDORecipes mocks base method.
func (m *MockDB) GetBDORecipes(skip, limit *int64) ([]data.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBDORecipes", skip, limit)
	ret0, _ := ret[0].([]data.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBDORecipes indicates an expected call of GetBDORecipes.
func (mr *MockDBMockRecorder) GetBDORecipes(skip, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBDORecipes", reflect.TypeOf((*MockDB)(nil).GetBDORecipes), skip, limit)
}

// GetDSPRecipes mocks base method.
func (m *MockDB) GetDSPRecipes(skip, limit *int64) ([]data.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDSPRecipes", skip, limit)
	ret0, _ := ret[0].([]data.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDSPRecipes indicates an expected call of GetDSPRecipes.
func (mr *MockDBMockRecorder) GetDSPRecipes(skip, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDSPRecipes", reflect.TypeOf((*MockDB)(nil).GetDSPRecipes), skip, limit)
}
