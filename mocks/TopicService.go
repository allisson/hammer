// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	hammer "github.com/allisson/hammer"
	mock "github.com/stretchr/testify/mock"
)

// TopicService is an autogenerated mock type for the TopicService type
type TopicService struct {
	mock.Mock
}

// Create provides a mock function with given fields: topic
func (_m *TopicService) Create(topic *hammer.Topic) error {
	ret := _m.Called(topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(*hammer.Topic) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *TopicService) Find(id string) (hammer.Topic, error) {
	ret := _m.Called(id)

	var r0 hammer.Topic
	if rf, ok := ret.Get(0).(func(string) hammer.Topic); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(hammer.Topic)
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
func (_m *TopicService) FindAll(findOptions hammer.FindOptions) ([]hammer.Topic, error) {
	ret := _m.Called(findOptions)

	var r0 []hammer.Topic
	if rf, ok := ret.Get(0).(func(hammer.FindOptions) []hammer.Topic); ok {
		r0 = rf(findOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]hammer.Topic)
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

// Update provides a mock function with given fields: topic
func (_m *TopicService) Update(topic *hammer.Topic) error {
	ret := _m.Called(topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(*hammer.Topic) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
