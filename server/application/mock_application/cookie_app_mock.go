// Code generated by MockGen. DO NOT EDIT.
// Source: application/cookie_app.go

// Package mock_application is a generated GoMock package.
package mock_application

import (
	http "net/http"
	entity "pinterest/domain/entity"
	__ "pinterest/services/auth/proto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCookieAppInterface is a mock of CookieAppInterface interface.
type MockCookieAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCookieAppInterfaceMockRecorder
}

// MockCookieAppInterfaceMockRecorder is the mock recorder for MockCookieAppInterface.
type MockCookieAppInterfaceMockRecorder struct {
	mock *MockCookieAppInterface
}

// NewMockCookieAppInterface creates a new mock instance.
func NewMockCookieAppInterface(ctrl *gomock.Controller) *MockCookieAppInterface {
	mock := &MockCookieAppInterface{ctrl: ctrl}
	mock.recorder = &MockCookieAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCookieAppInterface) EXPECT() *MockCookieAppInterfaceMockRecorder {
	return m.recorder
}

// AddCookieInfo mocks base method.
func (m *MockCookieAppInterface) AddCookieInfo(cookieInfo *entity.CookieInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCookieInfo", cookieInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCookieInfo indicates an expected call of AddCookieInfo.
func (mr *MockCookieAppInterfaceMockRecorder) AddCookieInfo(cookieInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCookieInfo", reflect.TypeOf((*MockCookieAppInterface)(nil).AddCookieInfo), cookieInfo)
}

// GenerateCookie mocks base method.
func (m *MockCookieAppInterface) GenerateCookie() (*http.Cookie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCookie")
	ret0, _ := ret[0].(*http.Cookie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCookie indicates an expected call of GenerateCookie.
func (mr *MockCookieAppInterfaceMockRecorder) GenerateCookie() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCookie", reflect.TypeOf((*MockCookieAppInterface)(nil).GenerateCookie))
}

// RemoveCookie mocks base method.
func (m *MockCookieAppInterface) RemoveCookie(arg0 *__.CookieInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCookie", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCookie indicates an expected call of RemoveCookie.
func (mr *MockCookieAppInterfaceMockRecorder) RemoveCookie(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCookie", reflect.TypeOf((*MockCookieAppInterface)(nil).RemoveCookie), arg0)
}

// SearchByUserID mocks base method.
func (m *MockCookieAppInterface) SearchByUserID(userID int) (*__.CookieInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByUserID", userID)
	ret0, _ := ret[0].(*__.CookieInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// SearchByUserID indicates an expected call of SearchByUserID.
func (mr *MockCookieAppInterfaceMockRecorder) SearchByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByUserID", reflect.TypeOf((*MockCookieAppInterface)(nil).SearchByUserID), userID)
}

// SearchByValue mocks base method.
func (m *MockCookieAppInterface) SearchByValue(sessionValue string) (*__.CookieInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByValue", sessionValue)
	ret0, _ := ret[0].(*__.CookieInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// SearchByValue indicates an expected call of SearchByValue.
func (mr *MockCookieAppInterfaceMockRecorder) SearchByValue(sessionValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByValue", reflect.TypeOf((*MockCookieAppInterface)(nil).SearchByValue), sessionValue)
}