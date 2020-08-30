package cmp

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual proto.Message) {
	if proto.Equal(expected, actual) {
		return
	}

	expectedFields := parse(expected)
	actualFields := parse(actual)

	err := matchFields(expectedFields, actualFields)
	if err != nil {
		t.Error(err)
	}
}
