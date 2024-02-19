// Code generated by mockery v2.14.0. DO NOT EDIT.

package memorizer

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockUtil is an autogenerated mock type for the Util type
type MockUtil[V interface{}] struct {
	mock.Mock
}

// Delete provides a mock function with given fields: k
func (_m *MockUtil[V]) Delete(k string) {
	_m.Called(k)
}

// Get provides a mock function with given fields: ctx, k
func (_m *MockUtil[V]) Get(ctx context.Context, k string) (V, error) {
	ret := _m.Called(ctx, k)

	var r0 V
	if rf, ok := ret.Get(0).(func(context.Context, string) V); ok {
		r0 = rf(ctx, k)
	} else {
		r0 = ret.Get(0).(V)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, k)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockUtil interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUtil creates a new instance of MockUtil. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUtil[V interface{}](t mockConstructorTestingTNewMockUtil) *MockUtil[V] {
	mock := &MockUtil[V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}