// Code generated by mockery v2.9.4. DO NOT EDIT.

package pheader

import mock "github.com/stretchr/testify/mock"

// mockMockConstructorTestingTNewMockPartnerRequest is an autogenerated mock type for the mockConstructorTestingTNewMockPartnerRequest type
type mockMockConstructorTestingTNewMockPartnerRequest struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: _a0
func (_m *mockMockConstructorTestingTNewMockPartnerRequest) Cleanup(_a0 func()) {
	_m.Called(_a0)
}

// Errorf provides a mock function with given fields: format, args
func (_m *mockMockConstructorTestingTNewMockPartnerRequest) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// FailNow provides a mock function with given fields:
func (_m *mockMockConstructorTestingTNewMockPartnerRequest) FailNow() {
	_m.Called()
}

// Logf provides a mock function with given fields: format, args
func (_m *mockMockConstructorTestingTNewMockPartnerRequest) Logf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}
