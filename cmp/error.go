package cmp

import "fmt"

type MatchError struct {
	err        error
	fiExpected []*fieldInfo
	fiActual   []*fieldInfo
}

func (m MatchError) Error() string {
	return fmt.Sprintf("%s\n++ expected\n%s\n\n-- actual\n%s", m.err, fieldsToString(m.fiExpected), fieldsToString(m.fiActual))
}

func newMatchError(err error, expected, actual []*fieldInfo) *MatchError {
	return &MatchError{
		err:        err,
		fiExpected: expected,
		fiActual:   actual,
	}
}

