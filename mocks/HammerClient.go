// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	api "github.com/allisson/hammer/api/v1"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// HammerClient is an autogenerated mock type for the HammerClient type
type HammerClient struct {
	mock.Mock
}

// CreateMessage provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) CreateMessage(ctx context.Context, in *api.CreateMessageRequest, opts ...grpc.CallOption) (*api.Message, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Message
	if rf, ok := ret.Get(0).(func(context.Context, *api.CreateMessageRequest, ...grpc.CallOption) *api.Message); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.CreateMessageRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSubscription provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) CreateSubscription(ctx context.Context, in *api.CreateSubscriptionRequest, opts ...grpc.CallOption) (*api.Subscription, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, *api.CreateSubscriptionRequest, ...grpc.CallOption) *api.Subscription); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.CreateSubscriptionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTopic provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) CreateTopic(ctx context.Context, in *api.CreateTopicRequest, opts ...grpc.CallOption) (*api.Topic, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Topic
	if rf, ok := ret.Get(0).(func(context.Context, *api.CreateTopicRequest, ...grpc.CallOption) *api.Topic); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.CreateTopicRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteMessage provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) DeleteMessage(ctx context.Context, in *api.DeleteMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *api.DeleteMessageRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.DeleteMessageRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSubscription provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) DeleteSubscription(ctx context.Context, in *api.DeleteSubscriptionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *api.DeleteSubscriptionRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.DeleteSubscriptionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTopic provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) DeleteTopic(ctx context.Context, in *api.DeleteTopicRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *api.DeleteTopicRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.DeleteTopicRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDelivery provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) GetDelivery(ctx context.Context, in *api.GetDeliveryRequest, opts ...grpc.CallOption) (*api.Delivery, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Delivery
	if rf, ok := ret.Get(0).(func(context.Context, *api.GetDeliveryRequest, ...grpc.CallOption) *api.Delivery); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.GetDeliveryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeliveryAttempt provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) GetDeliveryAttempt(ctx context.Context, in *api.GetDeliveryAttemptRequest, opts ...grpc.CallOption) (*api.DeliveryAttempt, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.DeliveryAttempt
	if rf, ok := ret.Get(0).(func(context.Context, *api.GetDeliveryAttemptRequest, ...grpc.CallOption) *api.DeliveryAttempt); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.DeliveryAttempt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.GetDeliveryAttemptRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessage provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) GetMessage(ctx context.Context, in *api.GetMessageRequest, opts ...grpc.CallOption) (*api.Message, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Message
	if rf, ok := ret.Get(0).(func(context.Context, *api.GetMessageRequest, ...grpc.CallOption) *api.Message); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.GetMessageRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscription provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) GetSubscription(ctx context.Context, in *api.GetSubscriptionRequest, opts ...grpc.CallOption) (*api.Subscription, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, *api.GetSubscriptionRequest, ...grpc.CallOption) *api.Subscription); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.GetSubscriptionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopic provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) GetTopic(ctx context.Context, in *api.GetTopicRequest, opts ...grpc.CallOption) (*api.Topic, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Topic
	if rf, ok := ret.Get(0).(func(context.Context, *api.GetTopicRequest, ...grpc.CallOption) *api.Topic); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.GetTopicRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDeliveries provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) ListDeliveries(ctx context.Context, in *api.ListDeliveriesRequest, opts ...grpc.CallOption) (*api.ListDeliveriesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.ListDeliveriesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *api.ListDeliveriesRequest, ...grpc.CallOption) *api.ListDeliveriesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ListDeliveriesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.ListDeliveriesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDeliveryAttempts provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) ListDeliveryAttempts(ctx context.Context, in *api.ListDeliveryAttemptsRequest, opts ...grpc.CallOption) (*api.ListDeliveryAttemptsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.ListDeliveryAttemptsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *api.ListDeliveryAttemptsRequest, ...grpc.CallOption) *api.ListDeliveryAttemptsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ListDeliveryAttemptsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.ListDeliveryAttemptsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListMessages provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) ListMessages(ctx context.Context, in *api.ListMessagesRequest, opts ...grpc.CallOption) (*api.ListMessagesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.ListMessagesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *api.ListMessagesRequest, ...grpc.CallOption) *api.ListMessagesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ListMessagesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.ListMessagesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSubscriptions provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) ListSubscriptions(ctx context.Context, in *api.ListSubscriptionsRequest, opts ...grpc.CallOption) (*api.ListSubscriptionsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.ListSubscriptionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *api.ListSubscriptionsRequest, ...grpc.CallOption) *api.ListSubscriptionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ListSubscriptionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.ListSubscriptionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTopics provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) ListTopics(ctx context.Context, in *api.ListTopicsRequest, opts ...grpc.CallOption) (*api.ListTopicsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.ListTopicsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *api.ListTopicsRequest, ...grpc.CallOption) *api.ListTopicsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ListTopicsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.ListTopicsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSubscription provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) UpdateSubscription(ctx context.Context, in *api.UpdateSubscriptionRequest, opts ...grpc.CallOption) (*api.Subscription, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, *api.UpdateSubscriptionRequest, ...grpc.CallOption) *api.Subscription); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.UpdateSubscriptionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTopic provides a mock function with given fields: ctx, in, opts
func (_m *HammerClient) UpdateTopic(ctx context.Context, in *api.UpdateTopicRequest, opts ...grpc.CallOption) (*api.Topic, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *api.Topic
	if rf, ok := ret.Get(0).(func(context.Context, *api.UpdateTopicRequest, ...grpc.CallOption) *api.Topic); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.UpdateTopicRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
