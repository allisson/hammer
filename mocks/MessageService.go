// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	hammer "github.com/allisson/hammer"
	mock "github.com/stretchr/testify/mock"
)

// MessageService is an autogenerated mock type for the MessageService type
type MessageService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, message
func (_m *MessageService) Create(ctx context.Context, message *hammer.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *hammer.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MessageService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, id
func (_m *MessageService) Find(ctx context.Context, id string) (*hammer.Message, error) {
	ret := _m.Called(ctx, id)

	var r0 *hammer.Message
	if rf, ok := ret.Get(0).(func(context.Context, string) *hammer.Message); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*hammer.Message)
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
func (_m *MessageService) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Message, error) {
	ret := _m.Called(ctx, findOptions)

	var r0 []*hammer.Message
	if rf, ok := ret.Get(0).(func(context.Context, hammer.FindOptions) []*hammer.Message); ok {
		r0 = rf(ctx, findOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*hammer.Message)
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
