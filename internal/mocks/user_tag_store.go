// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	model "article-tag/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserTagStore is an autogenerated mock type for the UserTagStore type
type UserTagStore struct {
	mock.Mock
}

type UserTagStore_Expecter struct {
	mock *mock.Mock
}

func (_m *UserTagStore) EXPECT() *UserTagStore_Expecter {
	return &UserTagStore_Expecter{mock: &_m.Mock}
}

// CreateTable provides a mock function with given fields: ctx
func (_m *UserTagStore) CreateTable(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserTagStore_CreateTable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTable'
type UserTagStore_CreateTable_Call struct {
	*mock.Call
}

// CreateTable is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserTagStore_Expecter) CreateTable(ctx interface{}) *UserTagStore_CreateTable_Call {
	return &UserTagStore_CreateTable_Call{Call: _e.mock.On("CreateTable", ctx)}
}

func (_c *UserTagStore_CreateTable_Call) Run(run func(ctx context.Context)) *UserTagStore_CreateTable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserTagStore_CreateTable_Call) Return(_a0 error) *UserTagStore_CreateTable_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserTagStore_CreateTable_Call) RunAndReturn(run func(context.Context) error) *UserTagStore_CreateTable_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, item
func (_m *UserTagStore) Delete(ctx context.Context, item *model.UserTag) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserTagStore_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserTagStore_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - item *model.UserTag
func (_e *UserTagStore_Expecter) Delete(ctx interface{}, item interface{}) *UserTagStore_Delete_Call {
	return &UserTagStore_Delete_Call{Call: _e.mock.On("Delete", ctx, item)}
}

func (_c *UserTagStore_Delete_Call) Run(run func(ctx context.Context, item *model.UserTag)) *UserTagStore_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserTag))
	})
	return _c
}

func (_c *UserTagStore_Delete_Call) Return(_a0 error) *UserTagStore_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserTagStore_Delete_Call) RunAndReturn(run func(context.Context, *model.UserTag) error) *UserTagStore_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DescribeTable provides a mock function with given fields: ctx
func (_m *UserTagStore) DescribeTable(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserTagStore_DescribeTable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DescribeTable'
type UserTagStore_DescribeTable_Call struct {
	*mock.Call
}

// DescribeTable is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserTagStore_Expecter) DescribeTable(ctx interface{}) *UserTagStore_DescribeTable_Call {
	return &UserTagStore_DescribeTable_Call{Call: _e.mock.On("DescribeTable", ctx)}
}

func (_c *UserTagStore_DescribeTable_Call) Run(run func(ctx context.Context)) *UserTagStore_DescribeTable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserTagStore_DescribeTable_Call) Return(_a0 error) *UserTagStore_DescribeTable_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserTagStore_DescribeTable_Call) RunAndReturn(run func(context.Context) error) *UserTagStore_DescribeTable_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, item
func (_m *UserTagStore) Get(ctx context.Context, item *model.UserTag) ([]string, error) {
	ret := _m.Called(ctx, item)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) ([]string, error)); ok {
		return rf(ctx, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) []string); ok {
		r0 = rf(ctx, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserTag) error); ok {
		r1 = rf(ctx, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserTagStore_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type UserTagStore_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - item *model.UserTag
func (_e *UserTagStore_Expecter) Get(ctx interface{}, item interface{}) *UserTagStore_Get_Call {
	return &UserTagStore_Get_Call{Call: _e.mock.On("Get", ctx, item)}
}

func (_c *UserTagStore_Get_Call) Run(run func(ctx context.Context, item *model.UserTag)) *UserTagStore_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserTag))
	})
	return _c
}

func (_c *UserTagStore_Get_Call) Return(_a0 []string, _a1 error) *UserTagStore_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserTagStore_Get_Call) RunAndReturn(run func(context.Context, *model.UserTag) ([]string, error)) *UserTagStore_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetPopularTags provides a mock function with given fields: ctx, item
func (_m *UserTagStore) GetPopularTags(ctx context.Context, item *model.UserTag) ([]string, error) {
	ret := _m.Called(ctx, item)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) ([]string, error)); ok {
		return rf(ctx, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) []string); ok {
		r0 = rf(ctx, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserTag) error); ok {
		r1 = rf(ctx, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserTagStore_GetPopularTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPopularTags'
type UserTagStore_GetPopularTags_Call struct {
	*mock.Call
}

// GetPopularTags is a helper method to define mock.On call
//   - ctx context.Context
//   - item *model.UserTag
func (_e *UserTagStore_Expecter) GetPopularTags(ctx interface{}, item interface{}) *UserTagStore_GetPopularTags_Call {
	return &UserTagStore_GetPopularTags_Call{Call: _e.mock.On("GetPopularTags", ctx, item)}
}

func (_c *UserTagStore_GetPopularTags_Call) Run(run func(ctx context.Context, item *model.UserTag)) *UserTagStore_GetPopularTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserTag))
	})
	return _c
}

func (_c *UserTagStore_GetPopularTags_Call) Return(_a0 []string, _a1 error) *UserTagStore_GetPopularTags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserTagStore_GetPopularTags_Call) RunAndReturn(run func(context.Context, *model.UserTag) ([]string, error)) *UserTagStore_GetPopularTags_Call {
	_c.Call.Return(run)
	return _c
}

// Store provides a mock function with given fields: ctx, item
func (_m *UserTagStore) Store(ctx context.Context, item *model.UserTag) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserTag) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserTagStore_Store_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Store'
type UserTagStore_Store_Call struct {
	*mock.Call
}

// Store is a helper method to define mock.On call
//   - ctx context.Context
//   - item *model.UserTag
func (_e *UserTagStore_Expecter) Store(ctx interface{}, item interface{}) *UserTagStore_Store_Call {
	return &UserTagStore_Store_Call{Call: _e.mock.On("Store", ctx, item)}
}

func (_c *UserTagStore_Store_Call) Run(run func(ctx context.Context, item *model.UserTag)) *UserTagStore_Store_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserTag))
	})
	return _c
}

func (_c *UserTagStore_Store_Call) Return(_a0 error) *UserTagStore_Store_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserTagStore_Store_Call) RunAndReturn(run func(context.Context, *model.UserTag) error) *UserTagStore_Store_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserTagStore creates a new instance of UserTagStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserTagStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserTagStore {
	mock := &UserTagStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
