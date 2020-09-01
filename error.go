package protocmp

import (
	"fmt"
	"strings"

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

func (m *matchErr) ValuesSwap() *matchErr {
	t := m.expected
	m.expected = m.actual
	m.actual = t
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

func (m *matchErr) Error() string {
	return fmt.Sprintf("%s: %s\n+ %v\n- %v", strings.Join(m.fieldKeys, "."), m.message, m.expected, m.actual)
}
