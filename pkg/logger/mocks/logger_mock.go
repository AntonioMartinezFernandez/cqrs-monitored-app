// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	slog "log/slog"
)

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: ctx, message, items
func (_m *Logger) Debug(ctx context.Context, message string, items ...slog.Attr) {
	_va := make([]any, len(items))
	for _i := range items {
		_va[_i] = items[_i]
	}
	var _ca []any
	_ca = append(_ca, ctx, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: ctx, message, items
func (_m *Logger) Error(ctx context.Context, message string, items ...slog.Attr) {
	_va := make([]any, len(items))
	for _i := range items {
		_va[_i] = items[_i]
	}
	var _ca []any
	_ca = append(_ca, ctx, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: ctx, message, items
func (_m *Logger) Info(ctx context.Context, message string, items ...slog.Attr) {
	_va := make([]any, len(items))
	for _i := range items {
		_va[_i] = items[_i]
	}
	var _ca []any
	_ca = append(_ca, ctx, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: ctx, message, items
func (_m *Logger) Warn(ctx context.Context, message string, items ...slog.Attr) {
	_va := make([]any, len(items))
	for _i := range items {
		_va[_i] = items[_i]
	}
	var _ca []any
	_ca = append(_ca, ctx, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// NewLogger creates a new instance of Logger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
