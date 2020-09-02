package protocmp

import (
	"testing"

	"github.com/nbaztec/protocmp/protos/sample"
)

func TestAssertEqual(t *testing.T) {
	actual := makeInput(func(v *sample.Outer) {})
	AssertEqual(t, makeInput(nil), actual)
}
