// Code generated by mockery v2.50.0. DO NOT EDIT.

package inframocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// UserHandlerHttpInterface is an autogenerated mock type for the UserHandlerHttpInterface type
type UserHandlerHttpInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) Create(c *gin.Context) {
	_m.Called(c)
}

// Delete provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) Delete(c *gin.Context) {
	_m.Called(c)
}

// Get provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) Get(c *gin.Context) {
	_m.Called(c)
}

// GetByFilterMany provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) GetByFilterMany(c *gin.Context) {
	_m.Called(c)
}

// GetByFilterOne provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) GetByFilterOne(c *gin.Context) {
	_m.Called(c)
}

// GetById provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) GetById(c *gin.Context) {
	_m.Called(c)
}

// Update provides a mock function with given fields: c
func (_m *UserHandlerHttpInterface) Update(c *gin.Context) {
	_m.Called(c)
}

// NewUserHandlerHttpInterface creates a new instance of UserHandlerHttpInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserHandlerHttpInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserHandlerHttpInterface {
	mock := &UserHandlerHttpInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
