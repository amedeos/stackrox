// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
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

// AckKeysIndexed mocks base method.
func (m *MockStore) AckKeysIndexed(ctx context.Context, keys ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AckKeysIndexed", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AckKeysIndexed indicates an expected call of AckKeysIndexed.
func (mr *MockStoreMockRecorder) AckKeysIndexed(ctx interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckKeysIndexed", reflect.TypeOf((*MockStore)(nil).AckKeysIndexed), varargs...)
}

// DeleteByQuery mocks base method.
func (m *MockStore) DeleteByQuery(ctx context.Context, query *v1.Query) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByQuery", ctx, query)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByQuery indicates an expected call of DeleteByQuery.
func (mr *MockStoreMockRecorder) DeleteByQuery(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByQuery", reflect.TypeOf((*MockStore)(nil).DeleteByQuery), ctx, query)
}

// DeleteMany mocks base method.
func (m *MockStore) DeleteMany(ctx context.Context, id []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMany", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMany indicates an expected call of DeleteMany.
func (mr *MockStoreMockRecorder) DeleteMany(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMany", reflect.TypeOf((*MockStore)(nil).DeleteMany), ctx, id)
}

// Get mocks base method.
func (m *MockStore) Get(ctx context.Context, id string) (*storage.ProcessIndicator, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.ProcessIndicator)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), ctx, id)
}

// GetByQuery mocks base method.
func (m *MockStore) GetByQuery(ctx context.Context, q *v1.Query) ([]*storage.ProcessIndicator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByQuery", ctx, q)
	ret0, _ := ret[0].([]*storage.ProcessIndicator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByQuery indicates an expected call of GetByQuery.
func (mr *MockStoreMockRecorder) GetByQuery(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByQuery", reflect.TypeOf((*MockStore)(nil).GetByQuery), ctx, q)
}

// GetKeysToIndex mocks base method.
func (m *MockStore) GetKeysToIndex(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeysToIndex", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeysToIndex indicates an expected call of GetKeysToIndex.
func (mr *MockStoreMockRecorder) GetKeysToIndex(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeysToIndex", reflect.TypeOf((*MockStore)(nil).GetKeysToIndex), ctx)
}

// GetMany mocks base method.
func (m *MockStore) GetMany(ctx context.Context, ids []string) ([]*storage.ProcessIndicator, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ctx, ids)
	ret0, _ := ret[0].([]*storage.ProcessIndicator)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMany indicates an expected call of GetMany.
func (mr *MockStoreMockRecorder) GetMany(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockStore)(nil).GetMany), ctx, ids)
}

// UpsertMany mocks base method.
func (m *MockStore) UpsertMany(arg0 context.Context, arg1 []*storage.ProcessIndicator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMany", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertMany indicates an expected call of UpsertMany.
func (mr *MockStoreMockRecorder) UpsertMany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMany", reflect.TypeOf((*MockStore)(nil).UpsertMany), arg0, arg1)
}

// Walk mocks base method.
func (m *MockStore) Walk(arg0 context.Context, arg1 func(*storage.ProcessIndicator) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Walk", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Walk indicates an expected call of Walk.
func (mr *MockStoreMockRecorder) Walk(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Walk", reflect.TypeOf((*MockStore)(nil).Walk), arg0, arg1)
}
