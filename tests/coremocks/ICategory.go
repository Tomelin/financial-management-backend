// Code generated by mockery v2.50.0. DO NOT EDIT.

package coremocks

import (
	mock "github.com/stretchr/testify/mock"
	entity "github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

// ICategory is an autogenerated mock type for the ICategory type
type ICategory struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *ICategory) Create(_a0 *entity.CategoryResponse) (*entity.CategoryResponse, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.CategoryResponse) (*entity.CategoryResponse, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*entity.CategoryResponse) *entity.CategoryResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.CategoryResponse) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *ICategory) Delete(id *string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with no fields
func (_m *ICategory) Get() ([]entity.CategoryResponse, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]entity.CategoryResponse, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []entity.CategoryResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByFilterMany provides a mock function with given fields: key, value
func (_m *ICategory) GetByFilterMany(key string, value *string) ([]entity.CategoryResponse, error) {
	ret := _m.Called(key, value)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterMany")
	}

	var r0 []entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *string) ([]entity.CategoryResponse, error)); ok {
		return rf(key, value)
	}
	if rf, ok := ret.Get(0).(func(string, *string) []entity.CategoryResponse); ok {
		r0 = rf(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *string) error); ok {
		r1 = rf(key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByFilterOne provides a mock function with given fields: key, value
func (_m *ICategory) GetByFilterOne(key string, value *string) (*entity.CategoryResponse, error) {
	ret := _m.Called(key, value)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterOne")
	}

	var r0 *entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *string) (*entity.CategoryResponse, error)); ok {
		return rf(key, value)
	}
	if rf, ok := ret.Get(0).(func(string, *string) *entity.CategoryResponse); ok {
		r0 = rf(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *string) error); ok {
		r1 = rf(key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *ICategory) GetById(id *string) (*entity.CategoryResponse, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 *entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*string) (*entity.CategoryResponse, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(*string) *entity.CategoryResponse); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: data
func (_m *ICategory) Update(data *entity.CategoryResponse) (*entity.CategoryResponse, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *entity.CategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.CategoryResponse) (*entity.CategoryResponse, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(*entity.CategoryResponse) *entity.CategoryResponse); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.CategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.CategoryResponse) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICategory creates a new instance of ICategory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICategory(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICategory {
	mock := &ICategory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
