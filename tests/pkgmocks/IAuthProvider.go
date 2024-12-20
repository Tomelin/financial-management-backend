// Code generated by mockery v2.50.0. DO NOT EDIT.

package pkgmocks

import (
	gin "github.com/gin-gonic/gin"
	authProvider "github.com/synera-br/financial-management/src/backend/pkg/authProvider"

	goth "github.com/markbates/goth"

	http "net/http"

	mock "github.com/stretchr/testify/mock"

	sessions "github.com/gin-gonic/contrib/sessions"
)

// IAuthProvider is an autogenerated mock type for the IAuthProvider type
type IAuthProvider struct {
	mock.Mock
}

// Callback provides a mock function with given fields: w, r
func (_m *IAuthProvider) Callback(w http.ResponseWriter, r *http.Request) (*goth.User, error) {
	ret := _m.Called(w, r)

	if len(ret) == 0 {
		panic("no return value specified for Callback")
	}

	var r0 *goth.User
	var r1 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) (*goth.User, error)); ok {
		return rf(w, r)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) *goth.User); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*goth.User)
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) error); ok {
		r1 = rf(w, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsLoggedIn provides a mock function with given fields: w, r
func (_m *IAuthProvider) IsLoggedIn(w http.ResponseWriter, r *http.Request) (*goth.User, error) {
	ret := _m.Called(w, r)

	if len(ret) == 0 {
		panic("no return value specified for IsLoggedIn")
	}

	var r0 *goth.User
	var r1 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) (*goth.User, error)); ok {
		return rf(w, r)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) *goth.User); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*goth.User)
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) error); ok {
		r1 = rf(w, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: w, r
func (_m *IAuthProvider) Login(w http.ResponseWriter, r *http.Request) (*authProvider.SessionStore, error) {
	ret := _m.Called(w, r)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 *authProvider.SessionStore
	var r1 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) (*authProvider.SessionStore, error)); ok {
		return rf(w, r)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) *authProvider.SessionStore); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authProvider.SessionStore)
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) error); ok {
		r1 = rf(w, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: c, w, r
func (_m *IAuthProvider) Logout(c *gin.Context, w http.ResponseWriter, r *http.Request) error {
	ret := _m.Called(c, w, r)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gin.Context, http.ResponseWriter, *http.Request) error); ok {
		r0 = rf(c, w, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with no fields
func (_m *IAuthProvider) Store() (*sessions.CookieStore, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Store")
	}

	var r0 *sessions.CookieStore
	var r1 error
	if rf, ok := ret.Get(0).(func() (*sessions.CookieStore, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *sessions.CookieStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.CookieStore)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAuthProvider creates a new instance of IAuthProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthProvider {
	mock := &IAuthProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
