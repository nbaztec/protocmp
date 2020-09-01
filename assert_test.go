package protocmp

import (
	"testing"

	"github.com/nbaztec/protocmp/protos/sample"
)

func TestAssertEqual(t *testing.T) {
	actual := patch(func(v *sample.Outer) {})
	AssertEqual(t, patchExpected(func(v *sample.Outer) {}), actual)
}
