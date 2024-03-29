// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	hammer "github.com/allisson/hammer"
	mock "github.com/stretchr/testify/mock"
)

// DeliveryService is an autogenerated mock type for the DeliveryService type
type DeliveryService struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, id
func (_m *DeliveryService) Find(ctx context.Context, id string) (*hammer.Delivery, error) {
	ret := _m.Called(ctx, id)

	var r0 *hammer.Delivery
	if rf, ok := ret.Get(0).(func(context.Context, string) *hammer.Delivery); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*hammer.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields: ctx, findOptions
func (_m *DeliveryService) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Delivery, error) {
	ret := _m.Called(ctx, findOptions)

	var r0 []*hammer.Delivery
	if rf, ok := ret.Get(0).(func(context.Context, hammer.FindOptions) []*hammer.Delivery); ok {
		r0 = rf(ctx, findOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*hammer.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, hammer.FindOptions) error); ok {
		r1 = rf(ctx, findOptions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
