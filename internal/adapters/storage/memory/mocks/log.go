// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	zapcore "go.uber.org/zap/zapcore"
)

// Log is an autogenerated mock type for the log type
type Log struct {
	mock.Mock
}

// Info provides a mock function with given fields: msg, fields
func (_m *Log) Info(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// NewLog creates a new instance of Log. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLog(t interface {
	mock.TestingT
	Cleanup(func())
}) *Log {
	mock := &Log{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
