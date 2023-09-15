// Code generated by mockery v2.23.2. DO NOT EDIT.

package fftmmocks

import mock "github.com/stretchr/testify/mock"

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Manager) Close() {
	_m.Called()
}

// Start provides a mock function with given fields:
func (_m *Manager) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewManager creates a new instance of Manager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewManager(t mockConstructorTestingTNewManager) *Manager {
	mock := &Manager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
