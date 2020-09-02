package protocmp

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/nbaztec/protocmp/protos/sample"
)

var now, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 00, time.UTC))

func TestAssertSuccess(t *testing.T) {
	check(
		t,
		makeInput(nil),
		makeInput(nil),
		nil,
	)
}

func TestAssertNil(t *testing.T) {
	check(
		t,
		makeInput(nil),
		nil,
		&DiffError{
			Field:    "Outer",
			Message:  "value mismatch",
			Expected: `<str_val:"foo" int_val:1 bool_val:true double_val:1.1 bytes_val:[1 2] repeated_type:[<id:"1"> <id:"2"> <nil>] map_type:map[A:<id:"AA"> B:<id:"BB"> C:<nil>] enum_type:NOT_OK oneof_string:"1" timestamp_type:<seconds:1598814300> duration_type:<seconds:1> any_type:<type_url:"mytype/v1" value:[5]> repeated_type_simple:[9 10 11] map_type_simple:map[A:20 B:30 C:40] nested_message:<inner:<id:"123">>>`,
			Actual:   `<nil>`,
		},
	)
}

func TestAssertString(t *testing.T) {
	check(
		t,
		makeInput(nil),
		makeInput(func(v *sample.Outer) {
			v.StrVal = "invalid"
		}),
		&DiffError{
			Field:    "str_val",
			Message:  "value mismatch",
			Expected: `"foo"`,
			Actual:   `"invalid"`,
		},
	)
}

func TestAssertInt(t *testing.T) {
	check(
		t,
		makeInput(nil),
		makeInput(func(v *sample.Outer) {
			v.IntVal = 42
		}),
		&DiffError{
			Field:    "int_val",
			Message:  "value mismatch",
			Expected: `1`,
			Actual:   `42`,
		},
	)
}

func TestAssertBool(t *testing.T) {
	check(
		t,
		makeInput(nil),
		makeInput(func(v *sample.Outer) {
			v.BoolVal = false
		}),
		&DiffError{
			Field:    "bool_val",
			Message:  "value mismatch",
			Expected: `true`,
			Actual:   `false`,
		},
	)
}

func TestAssertDouble(t *testing.T) {
	check(
		t,
		makeInput(nil),
		makeInput(func(v *sample.Outer) {
			v.DoubleVal = 42.1
		}),
		&DiffError{
			Field:    "double_val",
			Message:  "value mismatch",
			Expected: `1.1`,
			Actual:   `42.1`,
		},
	)
}

func TestAssertBytes(t *testing.T) {
	tests := []struct {
		name         string
		input        *sample.Outer
		diffExpected string
		diffActual   string
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.BytesVal = nil
			}),
			diffExpected: `[1 2]`,
			diffActual:   `[]`,
		},
		{
			name: "different length",
			input: makeInput(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6}
			}),
			diffExpected: `[1 2]`,
			diffActual:   `[6]`,
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6, 0x8}
			}),
			diffExpected: `[1 2]`,
			diffActual:   `[6 8]`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				&DiffError{
					Field:    "bytes_val",
					Message:  "value mismatch",
					Expected: tt.diffExpected,
					Actual:   tt.diffActual,
				},
			)
		})
	}
}

func TestAssertRepeated(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedType = nil
			}),
			diff: &DiffError{
				Field:    "repeated_type",
				Message:  "value mismatch",
				Expected: `[<id:"1"> <id:"2"> <nil>]`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different length",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "0"},
				}
			}),
			diff: &DiffError{
				Field:    "repeated_type",
				Message:  "length mismatch",
				Expected: `3`,
				Actual:   `1`,
			},
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					{Id: "3"},
					nil,
				}
			}),
			diff: &DiffError{
				Field:    "repeated_type.[1].id",
				Message:  "value mismatch",
				Expected: `"2"`,
				Actual:   `"3"`,
			},
		},
		{
			name: "different value - nil",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					nil,
					nil,
				}
			}),
			diff: &DiffError{
				Field:    "repeated_type.[1]",
				Message:  "value mismatch",
				Expected: `<id:"2">`,
				Actual:   `<nil>`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertRepeatedSimple(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedTypeSimple = nil
			}),
			diff: &DiffError{
				Field:    "repeated_type_simple",
				Message:  "value mismatch",
				Expected: `[9 10 11]`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different length",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{1}
			}),
			diff: &DiffError{
				Field:    "repeated_type_simple",
				Message:  "length mismatch",
				Expected: `3`,
				Actual:   `1`,
			},
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{9, 10, 1}
			}),
			diff: &DiffError{
				Field:    "repeated_type_simple.[2]",
				Message:  "value mismatch",
				Expected: `11`,
				Actual:   `1`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertMap(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.MapType = nil
			}),
			diff: &DiffError{
				Field:    "map_type",
				Message:  "value mismatch",
				Expected: `map[A:<id:"AA"> B:<id:"BB"> C:<nil>]`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different length",
			input: makeInput(func(v *sample.Outer) {
				v.MapType["X"] = nil
			}),
			diff: &DiffError{
				Field:    "map_type",
				Message:  "length mismatch",
				Expected: `3`,
				Actual:   `4`,
			},
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.MapType["B"].Id = "XYZ"
			}),
			diff: &DiffError{
				Field:    "map_type.[B].id",
				Message:  "value mismatch",
				Expected: `"BB"`,
				Actual:   `"XYZ"`,
			},
		},
		{
			name: "different value - nil",
			input: makeInput(func(v *sample.Outer) {
				v.MapType["B"] = nil
			}),
			diff: &DiffError{
				Field:    "map_type.[B]",
				Message:  "value mismatch",
				Expected: `<id:"BB">`,
				Actual:   `<nil>`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertMapSimple(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.MapTypeSimple = nil
			}),
			diff: &DiffError{
				Field:    "map_type_simple",
				Message:  "value mismatch",
				Expected: `map[A:20 B:30 C:40]`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different length",
			input: makeInput(func(v *sample.Outer) {
				v.MapTypeSimple["X"] = 0
			}),
			diff: &DiffError{
				Field:    "map_type_simple",
				Message:  "length mismatch",
				Expected: `3`,
				Actual:   `4`,
			},
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.MapTypeSimple["B"] = 99
			}),
			diff: &DiffError{
				Field:    "map_type_simple.[B]",
				Message:  "value mismatch",
				Expected: `30`,
				Actual:   `99`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertOneOf(t *testing.T) {
	tests := []struct {
		name     string
		expected *sample.Outer
		input    *sample.Outer
		diff     *DiffError
	}{
		{
			name: "nil value - simple",
			expected: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "XYZ"}
			}),
			input: makeInput(func(v *sample.Outer) {
				v.OneofType = nil
			}),
			diff: &DiffError{
				Field:    "oneof_string",
				Message:  "value mismatch",
				Expected: `"XYZ"`,
				Actual:   `""`,
			},
		},
		{
			name: "nil value - message",
			expected: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "XYZ"}}
			}),
			input: makeInput(func(v *sample.Outer) {
				v.OneofType = nil
			}),
			diff: &DiffError{
				Field:    "oneof_message",
				Message:  "value mismatch",
				Expected: `<id:"XYZ">`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "value - simple",
			expected: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "XYZ"}
			}),
			input: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "123"}
			}),
			diff: &DiffError{
				Field:    "oneof_string",
				Message:  "value mismatch",
				Expected: `"XYZ"`,
				Actual:   `"123"`,
			},
		},
		{
			name: "value - message",
			expected: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "XYZ"}}
			}),
			input: makeInput(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "123"}}
			}),
			diff: &DiffError{
				Field:    "oneof_message.id",
				Message:  "value mismatch",
				Expected: `"XYZ"`,
				Actual:   `"123"`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				tt.expected,
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertTimestamp(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.TimestampType = nil
			}),
			diff: &DiffError{
				Field:    "timestamp_type",
				Message:  "value mismatch",
				Expected: `<seconds:1598814300>`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different value - seconds",
			input: makeInput(func(v *sample.Outer) {
				v.TimestampType, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 10, 00, time.UTC))
			}),
			diff: &DiffError{
				Field:    "timestamp_type.seconds",
				Message:  "value mismatch",
				Expected: `1598814300`,
				Actual:   `1598814310`,
			},
		},
		{
			name: "different value - nanos",
			input: makeInput(func(v *sample.Outer) {
				v.TimestampType, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 10, time.UTC))
			}),
			diff: &DiffError{
				Field:    "timestamp_type.nanos",
				Message:  "value mismatch",
				Expected: `0`,
				Actual:   `10`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertDuration(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.DurationType = nil
			}),
			diff: &DiffError{
				Field:    "duration_type",
				Message:  "value mismatch",
				Expected: `<seconds:1>`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different value - seconds",
			input: makeInput(func(v *sample.Outer) {
				v.DurationType = ptypes.DurationProto(2 * time.Second)
			}),
			diff: &DiffError{
				Field:    "duration_type.seconds",
				Message:  "value mismatch",
				Expected: `1`,
				Actual:   `2`,
			},
		},
		{
			name: "different value - nanos",
			input: makeInput(func(v *sample.Outer) {
				v.DurationType = ptypes.DurationProto(1005 * time.Millisecond)
			}),
			diff: &DiffError{
				Field:    "duration_type.nanos",
				Message:  "value mismatch",
				Expected: `0`,
				Actual:   `5000000`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertAny(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.AnyType = nil
			}),
			diff: &DiffError{
				Field:    "any_type",
				Message:  "value mismatch",
				Expected: `<type_url:"mytype/v1" value:[5]>`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "different value",
			input: makeInput(func(v *sample.Outer) {
				v.AnyType.TypeUrl = "foo"
			}),
			diff: &DiffError{
				Field:    "any_type.type_url",
				Message:  "value mismatch",
				Expected: `"mytype/v1"`,
				Actual:   `"foo"`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func TestAssertNestedMessage(t *testing.T) {
	tests := []struct {
		name  string
		input *sample.Outer
		diff  *DiffError
	}{
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.NestedMessage = nil
			}),
			diff: &DiffError{
				Field:    "nested_message",
				Message:  "value mismatch",
				Expected: `<inner:<id:"123">>`,
				Actual:   `<nil>`,
			},
		},
		{
			name: "nil value",
			input: makeInput(func(v *sample.Outer) {
				v.NestedMessage.Inner.Id = "foo"
			}),
			diff: &DiffError{
				Field:    "nested_message.inner.id",
				Message:  "value mismatch",
				Expected: `"123"`,
				Actual:   `"foo"`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				makeInput(nil),
				tt.input,
				tt.diff,
			)
		})
	}
}

func makeInput(f func(v *sample.Outer)) *sample.Outer {
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
		EnumType:      sample.Outer_NOT_OK,
		OneofType:     &sample.Outer_OneofString{OneofString: "1"},
		TimestampType: now,
		DurationType:  ptypes.DurationProto(1 * time.Second),
		AnyType: &any.Any{
			TypeUrl: "mytype/v1",
			Value:   []byte{0x05},
		},
		RepeatedTypeSimple: []int32{9, 10, 11},
		MapTypeSimple: map[string]int32{
			"A": 20,
			"B": 30,
			"C": 40,
		},
		NestedMessage: &sample.Outer_NestedInner{Inner: &sample.Outer_NestedInner_Inner{Id: "123"}},
	}

	if f != nil {
		f(v)
	}

	return v
}

func check(t *testing.T, expected *sample.Outer, actual *sample.Outer, expectedErr *DiffError) {
	actualErr := Equal(expected, actual)
	if !reflect.DeepEqual(expectedErr, actualErr) {
		t.Errorf("mismatch err\n++ want:\n%s\n-- got:\n%s", expectedErr, actualErr)
		return
	}

	// check inverse
	if expectedErr != nil {
		x := expectedErr.Actual
		expectedErr.Actual = expectedErr.Expected
		expectedErr.Expected = x
	}

	actualErr = Equal(actual, expected)
	if !reflect.DeepEqual(expectedErr, actualErr) {
		t.Errorf("(inverse) mismatch err\n++ want:\n%s\n-- got:\n%s", expectedErr, actualErr)
	}
}
