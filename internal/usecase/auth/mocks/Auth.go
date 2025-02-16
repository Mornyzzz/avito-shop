// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	entity "avito-shop/internal/entity"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *Auth) Login(_a0 context.Context, _a1 entity.User) (string, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) (string, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuth creates a new instance of Auth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuth(t interface {
	mock.TestingT
	Cleanup(func())
}) *Auth {
	mock := &Auth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
