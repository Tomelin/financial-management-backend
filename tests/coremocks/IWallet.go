// Code generated by mockery v2.50.0. DO NOT EDIT.

package coremocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

// IWallet is an autogenerated mock type for the IWallet type
type IWallet struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, userId, wallet
func (_m *IWallet) Create(ctx context.Context, userId *string, wallet *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) Delete(ctx context.Context, userId *string, walletId *string) *entity.ModuleError {
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
func (_m *IWallet) Get(ctx context.Context, userId *string) ([]entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) GetByFilterMany(ctx context.Context, userId *string, filter []entity.QueryDB) ([]entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) GetByFilterOne(ctx context.Context, userId *string, filter []entity.QueryDB) (*entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) GetByID(ctx context.Context, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) GetWalletByIdAndUserID(ctx context.Context, userId *string, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {
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

// Update provides a mock function with given fields: ctx, userId, data
func (_m *IWallet) Update(ctx context.Context, userId *string, data *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {
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
func (_m *IWallet) UpdateBalance(ctx context.Context, walletID *string, balance *float64) (*entity.WalletResponse, *entity.ModuleError) {
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

// NewIWallet creates a new instance of IWallet. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIWallet(t interface {
	mock.TestingT
	Cleanup(func())
}) *IWallet {
	mock := &IWallet{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}