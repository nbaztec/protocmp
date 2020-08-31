package cmp

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/nbaztec/protocmp/protos/sample"
)

var now, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 00, time.UTC))
var expected = &sample.Outer{
	StrVal:    "foo",
	IntVal:    1,
	BoolVal:   true,
	DoubleVal: 1.1,
	BytesVal:  []byte{0x01, 0x02},
	RepeatedType: []*sample.Outer_Inner{
		{Id: "1"},
		{Id: "2"},
		nil,
	},
	MapType: map[string]*sample.Outer_Inner{
		"A": {Id: "AA"},
		"B": {Id: "BB"},
		"C": nil,
	},
	EnumType:            sample.Outer_NOT_OK,
	OneofType:           &sample.Outer_OneofInt{OneofInt: 1},
	LastUpdated:         now,
	LastUpdatedDuration: ptypes.DurationProto(1 * time.Second),
	Details: &any.Any{
		TypeUrl: "mytype/v1",
		Value:   []byte{0x05},
	},
	RepeatedTypeSimple: []int32{9, 10, 11},
}

func TestAssertEqual(t *testing.T) {
	actual := patch(func(v *sample.Outer) {})
	AssertEqual(t, expected, actual)
}

func checkForError(t *testing.T, actual *sample.Outer, expectedActualField, expectedExpectedField interface{}, expectedErr string) {
	actualErr := ""
	if err := Equal(expected, actual); err != nil {
		actualErr = err.Error()
	}

	if expectedErr != actualErr {
		t.Errorf("mismatch err\n++ want:\n%s\n-- got:\n%s", expectedErr, actualErr)
	}
}

func TestAssertString(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.StrVal = "invalid"
	})

	checkForError(
		t,
		actual,
		actual,
		expected,
		`str_val: value mismatch
+ foo
- invalid`,
	)
}

func TestAssertInt(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.IntVal = 42
	})

	checkForError(
		t,
		actual,
		actual,
		expected,
		`int_val: value mismatch
+ 1
- 42`,
	)
}

func TestAssertBool(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.BoolVal = false
	})

	checkForError(
		t,
		actual,
		actual,
		expected,
		`bool_val: value mismatch
+ true
- false`,
	)
}

func TestAssertDouble(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.DoubleVal = 42.1
	})

	checkForError(
		t,
		actual,
		actual,
		expected,
		`double_val: value mismatch
+ 1.1
- 42.1`,
	)
}

func TestAssertBytes(t *testing.T) {
	tests := []struct {
		name                string
		input               *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected   interface{}
		expectedErr         string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = nil
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- []`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.BytesVal
			},
			errFieldsExpected: expected.BytesVal,
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- [6]`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6, 0x8}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.BytesVal
			},
			errFieldsExpected: expected.BytesVal,
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- [6 8]`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkForError(
				t,
				tt.input,
				tt.errFieldsActualFunc(tt.input),
				tt.errFieldsExpected,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertRepeated(t *testing.T) {
	tests := []struct {
		name                string
		input               *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected   interface{}
		expectedErr         string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedType = nil
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr: `repeated_type: value mismatch
+ [<id:"1"> <id:"2"> <nil>]
- <nil>`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "0"},
				}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.RepeatedType
			},
			errFieldsExpected: expected.RepeatedType,
			expectedErr: `repeated_type: length mismatch
+ 3
- 1`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					{Id: "3"},
					nil,
				}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.RepeatedType[1]
			},
			errFieldsExpected: expected.RepeatedType[1],
			expectedErr: `repeated_type.[1].id: value mismatch
+ 2
- 3`,
		},
		{
			name: "different value - nil",
			input: patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					nil,
					nil,
				}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.RepeatedType
			},
			errFieldsExpected: expected.RepeatedType,
			expectedErr: `repeated_type.[1]: value mismatch
+ <id:"2">
- <nil>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkForError(
				t,
				tt.input,
				tt.errFieldsActualFunc(tt.input),
				tt.errFieldsExpected,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertRepeatedSimple(t *testing.T) {
	tests := []struct {
		name                string
		input               *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected   interface{}
		expectedErr         string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = nil
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr: `repeated_type_simple: value mismatch
+ [9 10 11]
- <nil>`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{1}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.RepeatedTypeSimple
			},
			errFieldsExpected: expected.RepeatedTypeSimple,
			expectedErr: `repeated_type_simple: length mismatch
+ 3
- 1`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{9, 10, 1}
			}),
			errFieldsActualFunc: func(v *sample.Outer) interface{} {
				return v.RepeatedTypeSimple
			},
			errFieldsExpected: expected.RepeatedTypeSimple,
			expectedErr: `repeated_type_simple.[2]: value mismatch
+ 11
- 1`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkForError(
				t,
				tt.input,
				tt.errFieldsActualFunc(tt.input),
				tt.errFieldsExpected,
				tt.expectedErr,
			)
		})
	}
}

func patch(f func(v *sample.Outer)) *sample.Outer {
	v := &sample.Outer{
		StrVal:    "foo",
		IntVal:    1,
		BoolVal:   true,
		DoubleVal: 1.1,
		BytesVal:  []byte{0x01, 0x02},
		RepeatedType: []*sample.Outer_Inner{
			{Id: "1"},
			{Id: "2"},
			nil,
		},
		MapType: map[string]*sample.Outer_Inner{
			"A": {Id: "AA"},
			"B": {Id: "BB"},
			"C": nil,
		},
		EnumType:            sample.Outer_NOT_OK,
		OneofType:           &sample.Outer_OneofInt{OneofInt: 1},
		LastUpdated:         now,
		LastUpdatedDuration: ptypes.DurationProto(1 * time.Second),
		Details: &any.Any{
			TypeUrl: "mytype/v1",
			Value:   []byte{0x05},
		},
		RepeatedTypeSimple: []int32{9, 10, 11},
	}

	f(v)

	return v
}
