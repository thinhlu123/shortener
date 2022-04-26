// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/thinhlu123/shortener/internal/models"
)

// MockLinkUsecases is a mock of LinkUsecases interface.
type MockLinkUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockLinkUsecasesMockRecorder
}

// MockLinkUsecasesMockRecorder is the mock recorder for MockLinkUsecases.
type MockLinkUsecasesMockRecorder struct {
	mock *MockLinkUsecases
}

// NewMockLinkUsecases creates a new mock instance.
func NewMockLinkUsecases(ctrl *gomock.Controller) *MockLinkUsecases {
	mock := &MockLinkUsecases{ctrl: ctrl}
	mock.recorder = &MockLinkUsecasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinkUsecases) EXPECT() *MockLinkUsecasesMockRecorder {
	return m.recorder
}

// CreateShortLink mocks base method.
func (m *MockLinkUsecases) CreateShortLink(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockLinkUsecasesMockRecorder) CreateShortLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink", reflect.TypeOf((*MockLinkUsecases)(nil).CreateShortLink), arg0, arg1)
}

// GetLink mocks base method.
func (m *MockLinkUsecases) GetLink(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLink", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLink indicates an expected call of GetLink.
func (mr *MockLinkUsecasesMockRecorder) GetLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLink", reflect.TypeOf((*MockLinkUsecases)(nil).GetLink), arg0, arg1)
}

// GetListLink mocks base method.
func (m *MockLinkUsecases) GetListLink(arg0 context.Context, arg1 models.Link) ([]models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListLink", arg0, arg1)
	ret0, _ := ret[0].([]models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListLink indicates an expected call of GetListLink.
func (mr *MockLinkUsecasesMockRecorder) GetListLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListLink", reflect.TypeOf((*MockLinkUsecases)(nil).GetListLink), arg0, arg1)
}
