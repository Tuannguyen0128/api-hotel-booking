// Code generated by mockery v2.14.0. DO NOT EDIT.

package encryptutil

import mock "github.com/stretchr/testify/mock"

// MockByteEncryptor is an autogenerated mock type for the ByteEncryptor type
type MockByteEncryptor struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: ciphertext
func (_m *MockByteEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	ret := _m.Called(ciphertext)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(ciphertext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(ciphertext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptWithUnmarshalling provides a mock function with given fields: ciphertext, v
func (_m *MockByteEncryptor) DecryptWithUnmarshalling(ciphertext []byte, v interface{}) error {
	ret := _m.Called(ciphertext, v)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, interface{}) error); ok {
		r0 = rf(ciphertext, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Encrypt provides a mock function with given fields: plaintext, isDeterministic
func (_m *MockByteEncryptor) Encrypt(plaintext []byte, isDeterministic bool) ([]byte, error) {
	ret := _m.Called(plaintext, isDeterministic)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, bool) []byte); ok {
		r0 = rf(plaintext, isDeterministic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, bool) error); ok {
		r1 = rf(plaintext, isDeterministic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptWithMarshalling provides a mock function with given fields: v, isDeterministic
func (_m *MockByteEncryptor) EncryptWithMarshalling(v interface{}, isDeterministic bool) ([]byte, error) {
	ret := _m.Called(v, isDeterministic)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(interface{}, bool) []byte); ok {
		r0 = rf(v, isDeterministic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, bool) error); ok {
		r1 = rf(v, isDeterministic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockByteEncryptor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockByteEncryptor creates a new instance of MockByteEncryptor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockByteEncryptor(t mockConstructorTestingTNewMockByteEncryptor) *MockByteEncryptor {
	mock := &MockByteEncryptor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
