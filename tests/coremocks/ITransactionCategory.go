// Code generated by mockery v2.50.0. DO NOT EDIT.

package coremocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

// ITransactionCategory is an autogenerated mock type for the ITransactionCategory type
type ITransactionCategory struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, email, category
func (_m *ITransactionCategory) Create(ctx context.Context, email *string, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, category)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, category)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.TransactionCategory) *entity.TransactionCategory); ok {
		r0 = rf(ctx, email, category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *entity.TransactionCategory) *entity.ModuleError); ok {
		r1 = rf(ctx, email, category)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, email, id
func (_m *ITransactionCategory) Delete(ctx context.Context, email *string, id *string) *entity.ModuleError {
	ret := _m.Called(ctx, email, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) *entity.ModuleError); ok {
		r0 = rf(ctx, email, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ModuleError)
		}
	}

	return r0
}

// Get provides a mock function with given fields: ctx, email, walletID
func (_m *ITransactionCategory) Get(ctx context.Context, email *string, walletID *string) ([]entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, walletID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) ([]entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, walletID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) []entity.TransactionCategory); ok {
		r0 = rf(ctx, email, walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *string) *entity.ModuleError); ok {
		r1 = rf(ctx, email, walletID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetByFilterMany provides a mock function with given fields: ctx, email, filter
func (_m *ITransactionCategory) GetByFilterMany(ctx context.Context, email *string, filter []entity.QueryDB) ([]entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterMany")
	}

	var r0 []entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) ([]entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) []entity.TransactionCategory); ok {
		r0 = rf(ctx, email, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, []entity.QueryDB) *entity.ModuleError); ok {
		r1 = rf(ctx, email, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetByFilterOne provides a mock function with given fields: ctx, email, filter
func (_m *ITransactionCategory) GetByFilterOne(ctx context.Context, email *string, filter []entity.QueryDB) (*entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterOne")
	}

	var r0 *entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) (*entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) *entity.TransactionCategory); ok {
		r0 = rf(ctx, email, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, []entity.QueryDB) *entity.ModuleError); ok {
		r1 = rf(ctx, email, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetById provides a mock function with given fields: ctx, email, id
func (_m *ITransactionCategory) GetById(ctx context.Context, email *string, id *string) (*entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 *entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) (*entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) *entity.TransactionCategory); ok {
		r0 = rf(ctx, email, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *string) *entity.ModuleError); ok {
		r1 = rf(ctx, email, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, email, category
func (_m *ITransactionCategory) Update(ctx context.Context, email *string, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	ret := _m.Called(ctx, email, category)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *entity.TransactionCategory
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError)); ok {
		return rf(ctx, email, category)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.TransactionCategory) *entity.TransactionCategory); ok {
		r0 = rf(ctx, email, category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.TransactionCategory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *entity.TransactionCategory) *entity.ModuleError); ok {
		r1 = rf(ctx, email, category)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// NewITransactionCategory creates a new instance of ITransactionCategory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewITransactionCategory(t interface {
	mock.TestingT
	Cleanup(func())
}) *ITransactionCategory {
	mock := &ITransactionCategory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
