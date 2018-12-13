// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/metrics/rules (interfaces: Store)

// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package rules is a generated GoMock package.
package rules

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockStore) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// ReadNamespaces mocks base method
func (m *MockStore) ReadNamespaces() (*Namespaces, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadNamespaces")
	ret0, _ := ret[0].(*Namespaces)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadNamespaces indicates an expected call of ReadNamespaces
func (mr *MockStoreMockRecorder) ReadNamespaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadNamespaces", reflect.TypeOf((*MockStore)(nil).ReadNamespaces))
}

// ReadRuleSet mocks base method
func (m *MockStore) ReadRuleSet(arg0 string) (RuleSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadRuleSet", arg0)
	ret0, _ := ret[0].(RuleSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadRuleSet indicates an expected call of ReadRuleSet
func (mr *MockStoreMockRecorder) ReadRuleSet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadRuleSet", reflect.TypeOf((*MockStore)(nil).ReadRuleSet), arg0)
}

// WriteAll mocks base method
func (m *MockStore) WriteAll(arg0 *Namespaces, arg1 MutableRuleSet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteAll indicates an expected call of WriteAll
func (mr *MockStoreMockRecorder) WriteAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAll", reflect.TypeOf((*MockStore)(nil).WriteAll), arg0, arg1)
}

// WriteRuleSet mocks base method
func (m *MockStore) WriteRuleSet(arg0 MutableRuleSet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteRuleSet", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteRuleSet indicates an expected call of WriteRuleSet
func (mr *MockStoreMockRecorder) WriteRuleSet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteRuleSet", reflect.TypeOf((*MockStore)(nil).WriteRuleSet), arg0)
}
