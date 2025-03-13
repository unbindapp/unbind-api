// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks_repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	repository "github.com/unbindapp/unbind-api/internal/repository"
)

// BaseRepositoryMock is an autogenerated mock type for the BaseRepositoryInterface type
type BaseRepositoryMock struct {
	mock.Mock
}

type BaseRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *BaseRepositoryMock) EXPECT() *BaseRepositoryMock_Expecter {
	return &BaseRepositoryMock_Expecter{mock: &_m.Mock}
}

// WithTx provides a mock function with given fields: ctx, fn
func (_m *BaseRepositoryMock) WithTx(ctx context.Context, fn func(repository.TxInterface) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(repository.TxInterface) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BaseRepositoryMock_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type BaseRepositoryMock_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(repository.TxInterface) error
func (_e *BaseRepositoryMock_Expecter) WithTx(ctx interface{}, fn interface{}) *BaseRepositoryMock_WithTx_Call {
	return &BaseRepositoryMock_WithTx_Call{Call: _e.mock.On("WithTx", ctx, fn)}
}

func (_c *BaseRepositoryMock_WithTx_Call) Run(run func(ctx context.Context, fn func(repository.TxInterface) error)) *BaseRepositoryMock_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(repository.TxInterface) error))
	})
	return _c
}

func (_c *BaseRepositoryMock_WithTx_Call) Return(_a0 error) *BaseRepositoryMock_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BaseRepositoryMock_WithTx_Call) RunAndReturn(run func(context.Context, func(repository.TxInterface) error) error) *BaseRepositoryMock_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewBaseRepositoryMock creates a new instance of BaseRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBaseRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *BaseRepositoryMock {
	mock := &BaseRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
