package protocmp

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/nbaztec/protocmp/protos/sample"
)

var now, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 00, time.UTC))

func TestAssertString(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.StrVal = "invalid"
	})

	checkStandard(
		t,
		actual,
		`str_val: value mismatch
+ "foo"
- "invalid"`,
	)
}

func TestAssertInt(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.IntVal = 42
	})

	checkStandard(
		t,
		actual,
		`int_val: value mismatch
+ 1
- 42`,
	)
}

func TestAssertBool(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.BoolVal = false
	})

	checkStandard(
		t,
		actual,
		`bool_val: value mismatch
+ true
- false`,
	)
}

func TestAssertDouble(t *testing.T) {
	actual := patch(func(v *sample.Outer) {
		v.DoubleVal = 42.1
	})

	checkStandard(
		t,
		actual,
		`double_val: value mismatch
+ 1.1
- 42.1`,
	)
}

func TestAssertBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = nil
			}),
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- []`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6}
			}),
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- [6]`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6, 0x8}
			}),
			expectedErr: `bytes_val: value mismatch
+ [1 2]
- [6 8]`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertRepeated(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedType = nil
			}),
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
			expectedErr: `repeated_type.[1].id: value mismatch
+ "2"
- "3"`,
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
			expectedErr: `repeated_type.[1]: value mismatch
+ <id:"2">
- <nil>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertRepeatedSimple(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = nil
			}),
			expectedErr: `repeated_type_simple: value mismatch
+ [9 10 11]
- <nil>`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{1}
			}),
			expectedErr: `repeated_type_simple: length mismatch
+ 3
- 1`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{9, 10, 1}
			}),
			expectedErr: `repeated_type_simple.[2]: value mismatch
+ 11
- 1`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertMap(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.MapType = nil
			}),
			expectedErr: `map_type: value mismatch
+ map[A:<id:"AA"> B:<id:"BB"> C:<nil>]
- <nil>`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.MapType["X"] = nil
			}),
			expectedErr: `map_type: length mismatch
+ 3
- 4`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.MapType["B"].Id = "XYZ"
			}),
			expectedErr: `map_type.[B].id: value mismatch
+ "BB"
- "XYZ"`,
		},
		{
			name: "different value - nil",
			input: patch(func(v *sample.Outer) {
				v.MapType["B"] = nil
			}),
			expectedErr: `map_type.[B]: value mismatch
+ <id:"BB">
- <nil>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertMapSimple(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.MapTypeSimple = nil
			}),
			expectedErr: `map_type_simple: value mismatch
+ map[A:20 B:30 C:40]
- <nil>`,
		},
		{
			name: "different length",
			input: patch(func(v *sample.Outer) {
				v.MapTypeSimple["X"] = 0
			}),
			expectedErr: `map_type_simple: length mismatch
+ 3
- 4`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.MapTypeSimple["B"] = 99
			}),
			expectedErr: `map_type_simple.[B]: value mismatch
+ 30
- 99`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertOneOf(t *testing.T) {
	tests := []struct {
		name        string
		expected    *sample.Outer
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value - simple",
			expected: patchExpected(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "XYZ"}
			}),
			input: patch(func(v *sample.Outer) {
				v.OneofType = nil
			}),
			expectedErr: `oneof_string: value mismatch
+ "XYZ"
- ""`,
		},
		{
			name: "nil value - message",
			expected: patchExpected(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "XYZ"}}
			}),
			input: patch(func(v *sample.Outer) {
				v.OneofType = nil
			}),
			expectedErr: `oneof_message: value mismatch
+ <id:"XYZ">
- <nil>`,
		},
		{
			name: "value - simple",
			expected: patchExpected(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "XYZ"}
			}),
			input: patch(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofString{OneofString: "123"}
			}),
			expectedErr: `oneof_string: value mismatch
+ "XYZ"
- "123"`,
		},
		{
			name: "value - message",
			expected: patchExpected(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "XYZ"}}
			}),
			input: patch(func(v *sample.Outer) {
				v.OneofType = &sample.Outer_OneofMessage{OneofMessage: &sample.Outer_Inner{Id: "123"}}
			}),
			expectedErr: `oneof_message.id: value mismatch
+ "XYZ"
- "123"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			check(
				t,
				tt.expected,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.TimestampType = nil
			}),
			expectedErr: `timestamp_type: value mismatch
+ <seconds:1598814300>
- <nil>`,
		},
		{
			name: "different value - seconds",
			input: patch(func(v *sample.Outer) {
				v.TimestampType, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 10, 00, time.UTC))
			}),
			expectedErr: `timestamp_type.seconds: value mismatch
+ 1598814300
- 1598814310`,
		},
		{
			name: "different value - nanos",
			input: patch(func(v *sample.Outer) {
				v.TimestampType, _ = ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 10, time.UTC))
			}),
			expectedErr: `timestamp_type.nanos: value mismatch
+ 0
- 10`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.DurationType = nil
			}),
			expectedErr: `duration_type: value mismatch
+ <seconds:1>
- <nil>`,
		},
		{
			name: "different value - seconds",
			input: patch(func(v *sample.Outer) {
				v.DurationType = ptypes.DurationProto(2 * time.Second)
			}),
			expectedErr: `duration_type.seconds: value mismatch
+ 1
- 2`,
		},
		{
			name: "different value - nanos",
			input: patch(func(v *sample.Outer) {
				v.DurationType = ptypes.DurationProto(1005 * time.Millisecond)
			}),
			expectedErr: `duration_type.nanos: value mismatch
+ 0
- 5000000`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertAny(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.AnyType = nil
			}),
			expectedErr: `any_type: value mismatch
+ <type_url:"mytype/v1" value:[5]>
- <nil>`,
		},
		{
			name: "different value",
			input: patch(func(v *sample.Outer) {
				v.AnyType.TypeUrl = "foo"
			}),
			expectedErr: `any_type.type_url: value mismatch
+ "mytype/v1"
- "foo"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func TestAssertNestedMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       *sample.Outer
		expectedErr string
	}{
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.NestedMessage = nil
			}),
			expectedErr: `nested_message: value mismatch
+ <inner:<id:"123">>
- <nil>`,
		},
		{
			name: "nil value",
			input: patch(func(v *sample.Outer) {
				v.NestedMessage.Inner.Id = "foo"
			}),
			expectedErr: `nested_message.inner.id: value mismatch
+ "123"
- "foo"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkStandard(
				t,
				tt.input,
				tt.expectedErr,
			)
		})
	}
}

func patchExpected(f func(v *sample.Outer)) *sample.Outer {
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

	f(v)

	return v
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

	f(v)

	return v
}

func checkStandard(t *testing.T, actual *sample.Outer, expectedErr string) {
	check(t, patchExpected(func(v *sample.Outer) {}), actual, expectedErr)
}

func check(t *testing.T, expected *sample.Outer, actual *sample.Outer, expectedErr string) {
	actualErr := ""
	if err := Equal(expected, actual); err != nil {
		actualErr = err.Error()
	}

	if expectedErr != actualErr {
		t.Errorf("mismatch err\n++ want:\n%s\n-- got:\n%s", expectedErr, actualErr)
	}
}
