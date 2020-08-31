package cmp

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual proto.Message) {
	err := Equal(expected, actual)
	if err != nil {
		t.Error(err)
	}
}
