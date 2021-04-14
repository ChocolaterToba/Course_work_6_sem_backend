// Code generated by MockGen. DO NOT EDIT.
// Source: application/pin_app.go

// Package mock_application is a generated GoMock package.
package mock_application

import (
	io "io"
	entity "pinterest/domain/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPinAppInterface is a mock of PinAppInterface interface.
type MockPinAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPinAppInterfaceMockRecorder
}

// MockPinAppInterfaceMockRecorder is the mock recorder for MockPinAppInterface.
type MockPinAppInterfaceMockRecorder struct {
	mock *MockPinAppInterface
}

// NewMockPinAppInterface creates a new mock instance.
func NewMockPinAppInterface(ctrl *gomock.Controller) *MockPinAppInterface {
	mock := &MockPinAppInterface{ctrl: ctrl}
	mock.recorder = &MockPinAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPinAppInterface) EXPECT() *MockPinAppInterfaceMockRecorder {
	return m.recorder
}

// AddPin mocks base method.
func (m *MockPinAppInterface) AddPin(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPin indicates an expected call of AddPin.
func (mr *MockPinAppInterfaceMockRecorder) AddPin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPin", reflect.TypeOf((*MockPinAppInterface)(nil).AddPin), arg0, arg1)
}

// CreatePin mocks base method.
func (m *MockPinAppInterface) CreatePin(arg0 int, arg1 *entity.Pin) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePin", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePin indicates an expected call of CreatePin.
func (mr *MockPinAppInterfaceMockRecorder) CreatePin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePin", reflect.TypeOf((*MockPinAppInterface)(nil).CreatePin), arg0, arg1)
}

// DeletePin mocks base method.
func (m *MockPinAppInterface) DeletePin(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePin indicates an expected call of DeletePin.
func (mr *MockPinAppInterfaceMockRecorder) DeletePin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePin", reflect.TypeOf((*MockPinAppInterface)(nil).DeletePin), arg0, arg1)
}

// GetLastUserPinID mocks base method.
func (m *MockPinAppInterface) GetLastUserPinID(arg0 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastUserPinID", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastUserPinID indicates an expected call of GetLastUserPinID.
func (mr *MockPinAppInterfaceMockRecorder) GetLastUserPinID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastUserPinID", reflect.TypeOf((*MockPinAppInterface)(nil).GetLastUserPinID), arg0)
}

// GetPin mocks base method.
func (m *MockPinAppInterface) GetPin(arg0 int) (*entity.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPin", arg0)
	ret0, _ := ret[0].(*entity.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPin indicates an expected call of GetPin.
func (mr *MockPinAppInterfaceMockRecorder) GetPin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPin", reflect.TypeOf((*MockPinAppInterface)(nil).GetPin), arg0)
}

// GetPins mocks base method.
func (m *MockPinAppInterface) GetPins(arg0 int) ([]entity.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPins", arg0)
	ret0, _ := ret[0].([]entity.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPins indicates an expected call of GetPins.
func (mr *MockPinAppInterfaceMockRecorder) GetPins(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPins", reflect.TypeOf((*MockPinAppInterface)(nil).GetPins), arg0)
}

// RemovePin mocks base method.
func (m *MockPinAppInterface) RemovePin(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePin indicates an expected call of RemovePin.
func (mr *MockPinAppInterfaceMockRecorder) RemovePin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePin", reflect.TypeOf((*MockPinAppInterface)(nil).RemovePin), arg0, arg1)
}

// SavePicture mocks base method.
func (m *MockPinAppInterface) SavePicture(arg0 *entity.Pin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePicture", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePicture indicates an expected call of SavePicture.
func (mr *MockPinAppInterfaceMockRecorder) SavePicture(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePicture", reflect.TypeOf((*MockPinAppInterface)(nil).SavePicture), arg0)
}

// SavePin mocks base method.
func (m *MockPinAppInterface) SavePin(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePin indicates an expected call of SavePin.
func (mr *MockPinAppInterfaceMockRecorder) SavePin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePin", reflect.TypeOf((*MockPinAppInterface)(nil).SavePin), arg0, arg1)
}

// UploadPicture mocks base method.
func (m *MockPinAppInterface) UploadPicture(arg0 int, arg1 io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadPicture", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadPicture indicates an expected call of UploadPicture.
func (mr *MockPinAppInterfaceMockRecorder) UploadPicture(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadPicture", reflect.TypeOf((*MockPinAppInterface)(nil).UploadPicture), arg0, arg1)
}
