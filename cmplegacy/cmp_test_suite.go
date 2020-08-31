package cmplegacy

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/nbaztec/protocmp/protos/sample"
	"reflect"
	"testing"
	"time"
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
	err := Equal(t, expected, actual)
	expectedFiActual := parseRecursive(expectedActualField, nil)
	expectedFiExpected := parseRecursive(expectedExpectedField, nil)
	if !reflect.DeepEqual(expectedFiActual, err.fiActual) {
		t.Errorf("mismatch actual fields\n++ want:\n%s\n-- got:\n%s", fieldsToString(expectedFiActual), fieldsToString(err.fiActual))
	}
	if !reflect.DeepEqual(expectedFiExpected, err.fiExpected) {
		t.Errorf("mismatch expected fields\n++ want:\n%s\n-- got:\n%s", fieldsToString(expectedFiExpected), fieldsToString(err.fiExpected))
	}
	if expectedErr != err.mainErr.Error() {
		t.Errorf("mismatch err\n++ want:\n%s\n-- got:\n%s", expectedErr, err.mainErr.Error())
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
		`field value mismatch: StrVal
+ foo
- invalid
`,
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
		`field value mismatch: IntVal
+ 1
- 42
`,
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
		`field value mismatch: BoolVal
+ true
- false
`,
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
		`field value mismatch: DoubleVal
+ 1.1
- 42.1
`,
	)
}

func TestAssertBytes(t *testing.T) {
	tests := []struct{
		name string
		input *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected interface{}
		expectedErr string
	} {
		{
			name:              "nil value",
			input:             patch(func(v *sample.Outer) {
				v.BytesVal = nil
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr:       `field value mismatch: BytesVal
+
[0]: 1
[1]: 2
- <nil>
`,
		},
		{
			name:              "different length",
			input:             patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.BytesVal
			},
			errFieldsExpected: expected.BytesVal,
			expectedErr:       `field length mismatch: BytesVal
+ 2
- 1
`,
		},
		{
			name:              "different value",
			input:             patch(func(v *sample.Outer) {
				v.BytesVal = []byte{0x6, 0x8}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.BytesVal
			},
			errFieldsExpected: expected.BytesVal,
			expectedErr:       `field value mismatch: BytesVal.[0]
+ 1
- 6
`,
		},
	}

	for _, tt:=range tests {
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
	tests := []struct{
		name string
		input *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected interface{}
		expectedErr string
	} {
		{
			name:              "nil value",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedType = nil
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr:       `field value mismatch: RepeatedType
+
[0]: 
 Id: 1
[1]: 
 Id: 2
[2]: <nil>
- <nil>
`,
		},
		{
			name:              "different length",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "0"},
				}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.RepeatedType
			},
			errFieldsExpected: expected.RepeatedType,
			expectedErr:       `field length mismatch: RepeatedType
+ 3
- 1
`,
		},
		{
			name:              "different value",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					{Id: "3"},
					nil,
				}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.RepeatedType[1]
			},
			errFieldsExpected: expected.RepeatedType[1],
			expectedErr:       `field value mismatch: RepeatedType.[1].Id
+ 2
- 3
`,
		},
		{
			name:              "different value - nil",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedType = []*sample.Outer_Inner{
					{Id: "1"},
					nil,
					nil,
				}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.RepeatedType
			},
			errFieldsExpected: expected.RepeatedType,
			expectedErr:       `field value mismatch: RepeatedType.[1]
+
Id: 2
- <nil>
`,
		},
	}

	for _, tt:=range tests {
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
	tests := []struct{
		name string
		input *sample.Outer
		errFieldsActualFunc func(*sample.Outer) interface{}
		errFieldsExpected interface{}
		expectedErr string
	} {
		{
			name:              "nil value",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = nil
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v
			},
			errFieldsExpected: expected,
			expectedErr:       `field value mismatch: RepeatedTypeSimple
+
[0]: 9
[1]: 10
[2]: 11
- <nil>
`,
		},
		{
			name:              "different length",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{1}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.RepeatedTypeSimple
			},
			errFieldsExpected: expected.RepeatedTypeSimple,
			expectedErr:       `field length mismatch: RepeatedTypeSimple
+ 3
- 1
`,
		},
		{
			name:              "different value",
			input:             patch(func(v *sample.Outer) {
				v.RepeatedTypeSimple = []int32{9, 10, 1}
			}),
			errFieldsActualFunc:   func(v *sample.Outer) interface{} {
				return v.RepeatedTypeSimple
			},
			errFieldsExpected: expected.RepeatedTypeSimple,
			expectedErr:       `field value mismatch: RepeatedTypeSimple.[2]
+ 11
- 1
`,
		},
	}

	for _, tt:=range tests {
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
