package protocmp

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/nbaztec/protocmp/protos/sample"
)

func TestAssertEqual(t *testing.T) {
	actual := makeInput(func(v *sample.Outer) {})
	AssertEqual(t, makeInput(nil), actual)
	AssertEqual(t, actual, makeInput(nil))
}

func TestAssertEqualFails(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AssertEqual(&testingT{}, nil, makeInput(nil))

	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = old

	expectedOutput := `TestXYZ: assert_test.go:23
        Outer: value mismatch
            + <nil>
            - <str_val:"foo" int_val:1 bool_val:true double_val:1.1 bytes_val:[1 2] repeated_type:[<id:"1"> <id:"2"> <nil>] map_type:map[A:<id:"AA"> B:<id:"BB"> C:<nil>] enum_type:NOT_OK oneof_string:"1" timestamp_type:<seconds:1598814300> duration_type:<seconds:1> any_type:<type_url:"mytype/v1" value:[5]> repeated_type_simple:[9 10 11] map_type_simple:map[A:20 B:30 C:40] nested_message:<inner:<id:"123">>>`

	actualOutput := strings.TrimSpace(string(out))
	if expectedOutput != actualOutput {
		t.Errorf("output mismatch want\n%s\ngot\n%s", expectedOutput, actualOutput)
	}
}

func TestAssertEqualFailsInverse(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AssertEqual(&testingT{}, makeInput(nil), nil)

	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = old

	expectedOutput := `TestXYZ: assert_test.go:45
        Outer: value mismatch
            + <str_val:"foo" int_val:1 bool_val:true double_val:1.1 bytes_val:[1 2] repeated_type:[<id:"1"> <id:"2"> <nil>] map_type:map[A:<id:"AA"> B:<id:"BB"> C:<nil>] enum_type:NOT_OK oneof_string:"1" timestamp_type:<seconds:1598814300> duration_type:<seconds:1> any_type:<type_url:"mytype/v1" value:[5]> repeated_type_simple:[9 10 11] map_type_simple:map[A:20 B:30 C:40] nested_message:<inner:<id:"123">>>
            - <nil>`

	actualOutput := strings.TrimSpace(string(out))
	if expectedOutput != actualOutput {
		t.Errorf("output mismatch want\n%s\ngot\n%s", expectedOutput, actualOutput)
	}
}

type testingT struct {
}

func (t *testingT) Name() string {
	return "TestXYZ"
}

func (t *testingT) Fail() {
}

