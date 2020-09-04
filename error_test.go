package protocmp

import (
	"reflect"
	"testing"
)

func TestMatchErrReturnsDiffError(t *testing.T) {
	err := &matchErr{
		fieldKeys: []string{"foo", "bar"},
		message:   "some message",
		expected:  "1",
		actual:    "2",
	}

	actual := err.Error()
	expected := DiffError{
		Field:    "foo.bar",
		Message:  "some message",
		Expected: "+ 1",
		Actual:   "- 2",
	}

	if reflect.DeepEqual(expected, actual) {
		t.Errorf("mismatch: want %+v, got %+v", expected, actual)
	}

}
