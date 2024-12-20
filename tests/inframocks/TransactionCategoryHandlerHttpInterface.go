// Code generated by mockery v2.50.0. DO NOT EDIT.

package inframocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// TransactionCategoryHandlerHttpInterface is an autogenerated mock type for the TransactionCategoryHandlerHttpInterface type
type TransactionCategoryHandlerHttpInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: c
func (_m *TransactionCategoryHandlerHttpInterface) Create(c *gin.Context) {
	_m.Called(c)
}

// Get provides a mock function with given fields: c
func (_m *TransactionCategoryHandlerHttpInterface) Get(c *gin.Context) {
	_m.Called(c)
}

// GetByFilterMany provides a mock function with given fields: c
func (_m *TransactionCategoryHandlerHttpInterface) GetByFilterMany(c *gin.Context) {
	_m.Called(c)
}

// GetByFilterOne provides a mock function with given fields: c
func (_m *TransactionCategoryHandlerHttpInterface) GetByFilterOne(c *gin.Context) {
	_m.Called(c)
}

// GetById provides a mock function with given fields: c
func (_m *TransactionCategoryHandlerHttpInterface) GetById(c *gin.Context) {
	_m.Called(c)
}

// NewTransactionCategoryHandlerHttpInterface creates a new instance of TransactionCategoryHandlerHttpInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionCategoryHandlerHttpInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionCategoryHandlerHttpInterface {
	mock := &TransactionCategoryHandlerHttpInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}