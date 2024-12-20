// Code generated by mockery v2.50.0. DO NOT EDIT.

package coremocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

// IWalletSvc is an autogenerated mock type for the IWalletSvc type
type IWalletSvc struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userId, wallet
func (_m *IWalletSvc) Create(ctx context.Context, userId *string, wallet *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId, wallet)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId, wallet)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.WalletResponse) *entity.WalletResponse); ok {
		r0 = rf(ctx, userId, wallet)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *entity.WalletResponse) *entity.ModuleError); ok {
		r1 = rf(ctx, userId, wallet)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, userId, walletId
func (_m *IWalletSvc) Delete(ctx context.Context, userId *string, walletId *string) *entity.ModuleError {
	ret := _m.Called(ctx, userId, walletId)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) *entity.ModuleError); ok {
		r0 = rf(ctx, userId, walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ModuleError)
		}
	}

	return r0
}

// Get provides a mock function with given fields: ctx, userId
func (_m *IWalletSvc) Get(ctx context.Context, userId *string) ([]entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string) ([]entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) []entity.WalletResponse); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) *entity.ModuleError); ok {
		r1 = rf(ctx, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetByFilterMany provides a mock function with given fields: ctx, userId, filter
func (_m *IWalletSvc) GetByFilterMany(ctx context.Context, userId *string, filter []entity.QueryDB) ([]entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterMany")
	}

	var r0 []entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) ([]entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) []entity.WalletResponse); ok {
		r0 = rf(ctx, userId, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, []entity.QueryDB) *entity.ModuleError); ok {
		r1 = rf(ctx, userId, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetByFilterOne provides a mock function with given fields: ctx, userId, filter
func (_m *IWalletSvc) GetByFilterOne(ctx context.Context, userId *string, filter []entity.QueryDB) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetByFilterOne")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, []entity.QueryDB) *entity.WalletResponse); ok {
		r0 = rf(ctx, userId, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, []entity.QueryDB) *entity.ModuleError); ok {
		r1 = rf(ctx, userId, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, walletId
func (_m *IWalletSvc) GetByID(ctx context.Context, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, walletId)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, walletId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *entity.WalletResponse); ok {
		r0 = rf(ctx, walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) *entity.ModuleError); ok {
		r1 = rf(ctx, walletId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// GetWalletByIdAndUserID provides a mock function with given fields: ctx, userId, walletId
func (_m *IWalletSvc) GetWalletByIdAndUserID(ctx context.Context, userId *string, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId, walletId)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletByIdAndUserID")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId, walletId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *string) *entity.WalletResponse); ok {
		r0 = rf(ctx, userId, walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *string) *entity.ModuleError); ok {
		r1 = rf(ctx, userId, walletId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// PlanValidate provides a mock function with given fields: id
func (_m *IWalletSvc) PlanValidate(id *string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for PlanValidate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetBalance provides a mock function with given fields: id, balance
func (_m *IWalletSvc) SetBalance(id *string, balance *float64) error {
	ret := _m.Called(id, balance)

	if len(ret) == 0 {
		panic("no return value specified for SetBalance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*string, *float64) error); ok {
		r0 = rf(id, balance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, userId, data
func (_m *IWalletSvc) Update(ctx context.Context, userId *string, data *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, userId, data)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, userId, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *entity.WalletResponse) *entity.WalletResponse); ok {
		r0 = rf(ctx, userId, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *entity.WalletResponse) *entity.ModuleError); ok {
		r1 = rf(ctx, userId, data)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// UpdateBalance provides a mock function with given fields: ctx, walletID, balance
func (_m *IWalletSvc) UpdateBalance(ctx context.Context, walletID *string, balance *float64) (*entity.WalletResponse, *entity.ModuleError) {
	ret := _m.Called(ctx, walletID, balance)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBalance")
	}

	var r0 *entity.WalletResponse
	var r1 *entity.ModuleError
	if rf, ok := ret.Get(0).(func(context.Context, *string, *float64) (*entity.WalletResponse, *entity.ModuleError)); ok {
		return rf(ctx, walletID, balance)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *float64) *entity.WalletResponse); ok {
		r0 = rf(ctx, walletID, balance)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.WalletResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *float64) *entity.ModuleError); ok {
		r1 = rf(ctx, walletID, balance)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*entity.ModuleError)
		}
	}

	return r0, r1
}

// NewIWalletSvc creates a new instance of IWalletSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIWalletSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *IWalletSvc {
	mock := &IWalletSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
