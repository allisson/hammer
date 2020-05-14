// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	hammer "github.com/allisson/hammer"
	mock "github.com/stretchr/testify/mock"
)

// SubscriptionService is an autogenerated mock type for the SubscriptionService type
type SubscriptionService struct {
	mock.Mock
}

// Create provides a mock function with given fields: subscription
func (_m *SubscriptionService) Create(subscription *hammer.Subscription) error {
	ret := _m.Called(subscription)

	var r0 error
	if rf, ok := ret.Get(0).(func(*hammer.Subscription) error); ok {
		r0 = rf(subscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *SubscriptionService) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *SubscriptionService) Find(id string) (hammer.Subscription, error) {
	ret := _m.Called(id)

	var r0 hammer.Subscription
	if rf, ok := ret.Get(0).(func(string) hammer.Subscription); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(hammer.Subscription)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields: findOptions
func (_m *SubscriptionService) FindAll(findOptions hammer.FindOptions) ([]hammer.Subscription, error) {
	ret := _m.Called(findOptions)

	var r0 []hammer.Subscription
	if rf, ok := ret.Get(0).(func(hammer.FindOptions) []hammer.Subscription); ok {
		r0 = rf(findOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]hammer.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(hammer.FindOptions) error); ok {
		r1 = rf(findOptions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: subscription
func (_m *SubscriptionService) Update(subscription *hammer.Subscription) error {
	ret := _m.Called(subscription)

	var r0 error
	if rf, ok := ret.Get(0).(func(*hammer.Subscription) error); ok {
		r0 = rf(subscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
