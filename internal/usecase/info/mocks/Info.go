// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	entity "avito-shop/internal/entity"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Info is an autogenerated mock type for the Info type
type Info struct {
	mock.Mock
}

// GetInfo provides a mock function with given fields: ctx, username
func (_m *Info) GetInfo(ctx context.Context, username string) (*entity.Info, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetInfo")
	}

	var r0 *entity.Info
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.Info, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Info); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Info)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewInfo creates a new instance of Info. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInfo(t interface {
	mock.TestingT
	Cleanup(func())
}) *Info {
	mock := &Info{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
