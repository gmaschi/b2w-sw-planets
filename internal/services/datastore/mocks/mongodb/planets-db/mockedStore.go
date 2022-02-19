// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db (interfaces: Store)

// Package mockedstore is a generated GoMock package.
package mockedstore

import (
	context "context"
	reflect "reflect"

	planetsdb "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreatePlanet mocks base method.
func (m *MockStore) CreatePlanet(arg0 context.Context, arg1 planetsdb.CreatePlanetParams) (planetsdb.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlanet", arg0, arg1)
	ret0, _ := ret[0].(planetsdb.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlanet indicates an expected call of CreatePlanet.
func (mr *MockStoreMockRecorder) CreatePlanet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlanet", reflect.TypeOf((*MockStore)(nil).CreatePlanet), arg0, arg1)
}

// DeletePlanet mocks base method.
func (m *MockStore) DeletePlanet(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlanet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlanet indicates an expected call of DeletePlanet.
func (mr *MockStoreMockRecorder) DeletePlanet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlanet", reflect.TypeOf((*MockStore)(nil).DeletePlanet), arg0, arg1)
}

// GetPlanet mocks base method.
func (m *MockStore) GetPlanet(arg0 context.Context, arg1 string) (planetsdb.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanet", arg0, arg1)
	ret0, _ := ret[0].(planetsdb.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanet indicates an expected call of GetPlanet.
func (mr *MockStoreMockRecorder) GetPlanet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanet", reflect.TypeOf((*MockStore)(nil).GetPlanet), arg0, arg1)
}

// ListPlanets mocks base method.
func (m *MockStore) ListPlanets(arg0 context.Context, arg1 planetsdb.ListPlanetParams) ([]planetsdb.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPlanets", arg0, arg1)
	ret0, _ := ret[0].([]planetsdb.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPlanets indicates an expected call of ListPlanets.
func (mr *MockStoreMockRecorder) ListPlanets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPlanets", reflect.TypeOf((*MockStore)(nil).ListPlanets), arg0, arg1)
}
