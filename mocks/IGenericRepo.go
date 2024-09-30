// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IGenericRepo is an autogenerated mock type for the IGenericRepo type
type IGenericRepo[T interface{}, X interface{ string | uint }] struct {
	mock.Mock
}

type IGenericRepo_Expecter[T interface{}, X interface{ string | uint }] struct {
	mock *mock.Mock
}

func (_m *IGenericRepo[T, X]) EXPECT() *IGenericRepo_Expecter[T, X] {
	return &IGenericRepo_Expecter[T, X]{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *IGenericRepo[T, X]) Create(_a0 context.Context, _a1 T) (T, error) {
	ret := _m.Called(_a0, _a1)

	var r0 T
	if rf, ok := ret.Get(0).(func(context.Context, T) T); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(T)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, T) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IGenericRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type IGenericRepo_Create_Call[T interface{}, X interface{ string | uint }] struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 T
func (_e *IGenericRepo_Expecter[T, X]) Create(_a0 interface{}, _a1 interface{}) *IGenericRepo_Create_Call[T, X] {
	return &IGenericRepo_Create_Call[T, X]{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *IGenericRepo_Create_Call[T, X]) Run(run func(_a0 context.Context, _a1 T)) *IGenericRepo_Create_Call[T, X] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(T))
	})
	return _c
}

func (_c *IGenericRepo_Create_Call[T, X]) Return(_a0 T, _a1 error) *IGenericRepo_Create_Call[T, X] {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *IGenericRepo[T, X]) Delete(_a0 context.Context, _a1 X, _a2 bool) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, X, bool) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IGenericRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type IGenericRepo_Delete_Call[T interface{}, X interface{ string | uint }] struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 X
//   - _a2 bool
func (_e *IGenericRepo_Expecter[T, X]) Delete(_a0 interface{}, _a1 interface{}, _a2 interface{}) *IGenericRepo_Delete_Call[T, X] {
	return &IGenericRepo_Delete_Call[T, X]{Call: _e.mock.On("Delete", _a0, _a1, _a2)}
}

func (_c *IGenericRepo_Delete_Call[T, X]) Run(run func(_a0 context.Context, _a1 X, _a2 bool)) *IGenericRepo_Delete_Call[T, X] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(X), args[2].(bool))
	})
	return _c
}

func (_c *IGenericRepo_Delete_Call[T, X]) Return(_a0 error) *IGenericRepo_Delete_Call[T, X] {
	_c.Call.Return(_a0)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1, _a2
func (_m *IGenericRepo[T, X]) Get(_a0 context.Context, _a1 X, _a2 string) (*T, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, X, string) *T); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, X, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IGenericRepo_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type IGenericRepo_Get_Call[T interface{}, X interface{ string | uint }] struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 X
//   - _a2 string
func (_e *IGenericRepo_Expecter[T, X]) Get(_a0 interface{}, _a1 interface{}, _a2 interface{}) *IGenericRepo_Get_Call[T, X] {
	return &IGenericRepo_Get_Call[T, X]{Call: _e.mock.On("Get", _a0, _a1, _a2)}
}

func (_c *IGenericRepo_Get_Call[T, X]) Run(run func(_a0 context.Context, _a1 X, _a2 string)) *IGenericRepo_Get_Call[T, X] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(X), args[2].(string))
	})
	return _c
}

func (_c *IGenericRepo_Get_Call[T, X]) Return(_a0 *T, _a1 error) *IGenericRepo_Get_Call[T, X] {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetAll provides a mock function with given fields: _a0
func (_m *IGenericRepo[T, X]) GetAll(_a0 context.Context) ([]*T, error) {
	ret := _m.Called(_a0)

	var r0 []*T
	if rf, ok := ret.Get(0).(func(context.Context) []*T); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*T)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IGenericRepo_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type IGenericRepo_GetAll_Call[T interface{}, X interface{ string | uint }] struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *IGenericRepo_Expecter[T, X]) GetAll(_a0 interface{}) *IGenericRepo_GetAll_Call[T, X] {
	return &IGenericRepo_GetAll_Call[T, X]{Call: _e.mock.On("GetAll", _a0)}
}

func (_c *IGenericRepo_GetAll_Call[T, X]) Run(run func(_a0 context.Context)) *IGenericRepo_GetAll_Call[T, X] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *IGenericRepo_GetAll_Call[T, X]) Return(_a0 []*T, _a1 error) *IGenericRepo_GetAll_Call[T, X] {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1, _a2
func (_m *IGenericRepo[T, X]) Update(_a0 context.Context, _a1 X, _a2 T) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, X, T) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IGenericRepo_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type IGenericRepo_Update_Call[T interface{}, X interface{ string | uint }] struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 X
//   - _a2 T
func (_e *IGenericRepo_Expecter[T, X]) Update(_a0 interface{}, _a1 interface{}, _a2 interface{}) *IGenericRepo_Update_Call[T, X] {
	return &IGenericRepo_Update_Call[T, X]{Call: _e.mock.On("Update", _a0, _a1, _a2)}
}

func (_c *IGenericRepo_Update_Call[T, X]) Run(run func(_a0 context.Context, _a1 X, _a2 T)) *IGenericRepo_Update_Call[T, X] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(X), args[2].(T))
	})
	return _c
}

func (_c *IGenericRepo_Update_Call[T, X]) Return(_a0 error) *IGenericRepo_Update_Call[T, X] {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewIGenericRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewIGenericRepo creates a new instance of IGenericRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIGenericRepo[T interface{}, X interface{ string | uint }](t mockConstructorTestingTNewIGenericRepo) *IGenericRepo[T, X] {
	mock := &IGenericRepo[T, X]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
