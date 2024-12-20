// Code generated by mockery v2.50.0. DO NOT EDIT.

package pkgmocks

import (
	context "context"

	firestore "cloud.google.com/go/firestore"

	mock "github.com/stretchr/testify/mock"
)

// FirebaseDatabaseInterface is an autogenerated mock type for the FirebaseDatabaseInterface type
type FirebaseDatabaseInterface struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *FirebaseDatabaseInterface) Close() {
	_m.Called()
}

// Collection provides a mock function with given fields: name
func (_m *FirebaseDatabaseInterface) Collection(name string) *firestore.CollectionRef {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Collection")
	}

	var r0 *firestore.CollectionRef
	if rf, ok := ret.Get(0).(func(string) *firestore.CollectionRef); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*firestore.CollectionRef)
		}
	}

	return r0
}

// Documents provides a mock function with given fields: ctx, name
func (_m *FirebaseDatabaseInterface) Documents(ctx context.Context, name string) *firestore.DocumentIterator {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for Documents")
	}

	var r0 *firestore.DocumentIterator
	if rf, ok := ret.Get(0).(func(context.Context, string) *firestore.DocumentIterator); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*firestore.DocumentIterator)
		}
	}

	return r0
}

// NewFirebaseDatabaseInterface creates a new instance of FirebaseDatabaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFirebaseDatabaseInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *FirebaseDatabaseInterface {
	mock := &FirebaseDatabaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
