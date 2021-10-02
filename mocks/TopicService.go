// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	hammer "github.com/allisson/hammer"
	mock "github.com/stretchr/testify/mock"
)

// TopicService is an autogenerated mock type for the TopicService type
type TopicService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, topic
func (_m *TopicService) Create(ctx context.Context, topic *hammer.Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *hammer.Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TopicService) Delete(ctx context.Context, id string) error {
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
func (_m *TopicService) Find(ctx context.Context, id string) (*hammer.Topic, error) {
	ret := _m.Called(ctx, id)

	var r0 *hammer.Topic
	if rf, ok := ret.Get(0).(func(context.Context, string) *hammer.Topic); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*hammer.Topic)
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
func (_m *TopicService) FindAll(ctx context.Context, findOptions hammer.FindOptions) ([]*hammer.Topic, error) {
	ret := _m.Called(ctx, findOptions)

	var r0 []*hammer.Topic
	if rf, ok := ret.Get(0).(func(context.Context, hammer.FindOptions) []*hammer.Topic); ok {
		r0 = rf(ctx, findOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*hammer.Topic)
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

// Update provides a mock function with given fields: ctx, topic
func (_m *TopicService) Update(ctx context.Context, topic *hammer.Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *hammer.Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
