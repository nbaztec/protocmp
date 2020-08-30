package cmp

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func Equal(t *testing.T, expected, actual proto.Message) *MatchError {
	if proto.Equal(expected, actual) {
		return nil
	}

	actualFields := parse(actual)
	expectedFields := parse(expected)

	return matchFields(expectedFields, actualFields)
}

func AssertEqual(t *testing.T, expected, actual proto.Message) {
	err := Equal(t, expected, actual)
	if err != nil {
		t.Error(err)
	}
}
