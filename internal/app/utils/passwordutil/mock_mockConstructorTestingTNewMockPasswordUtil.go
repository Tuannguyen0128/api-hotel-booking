// Code generated by mockery v2.9.4. DO NOT EDIT.

package passwordutil

import mock "github.com/stretchr/testify/mock"

// mockMockConstructorTestingTNewMockPasswordUtil is an autogenerated mock type for the mockConstructorTestingTNewMockPasswordUtil type
type mockMockConstructorTestingTNewMockPasswordUtil struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: _a0
func (_m *mockMockConstructorTestingTNewMockPasswordUtil) Cleanup(_a0 func()) {
	_m.Called(_a0)
}

// Errorf provides a mock function with given fields: format, args
func (_m *mockMockConstructorTestingTNewMockPasswordUtil) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// FailNow provides a mock function with given fields:
func (_m *mockMockConstructorTestingTNewMockPasswordUtil) FailNow() {
	_m.Called()
}

// Logf provides a mock function with given fields: format, args
func (_m *mockMockConstructorTestingTNewMockPasswordUtil) Logf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}