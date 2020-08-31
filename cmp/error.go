package cmp

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type matchErr struct {
	fieldKeys []string
	message   string
	expected  interface{}
	actual    interface{}
}

func newMatchError(message string) *matchErr {
	return &matchErr{message: message}
}

func (m *matchErr) Field(k protoreflect.Name) *matchErr {
	m.fieldKeys = append([]string{string(k)}, m.fieldKeys...)
	return m
}

func (m *matchErr) Values(expected, actual interface{}) *matchErr {
	m.expected = expected
	m.actual = actual
	return m
}

func (m *matchErr) ValueActual(actual interface{}) *matchErr {
	m.actual = actual
	return m
}

func (m *matchErr) ValueExpected(expected interface{}) *matchErr {
	m.expected = expected
	return m
}

