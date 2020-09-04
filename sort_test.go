package protocmp

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"testing"
)

func TestSortedMapRange(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		inputKind protoreflect.Kind
		expected  interface{}
	}{
		{
			name:      "string keys",
			inputKind: protoreflect.StringKind,
			input:     map[string]int32{"z": 1, "a": 2, "c": 3},
			expected:  []interface{}{"a", "c", "z"},
		},
		{
			name:      "bool keys",
			inputKind: protoreflect.BoolKind,
			input:     map[bool]int32{false: 1, true: 2},
			expected:  []interface{}{false, true},
		},
		{
			name:      "int32 keys",
			inputKind: protoreflect.Int32Kind,
			input:     map[int32]int32{10: 1, 30: 2, 20: 3, -5: 4},
			expected:  []interface{}{int32(-5), int32(10), int32(20), int32(30)},
		},
		{
			name:      "uint32 keys",
			inputKind: protoreflect.Uint32Kind,
			input:     map[uint32]int32{10: 1, 30: 2, 20: 3},
			expected:  []interface{}{uint32(10), uint32(20), uint32(30)},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := newMapT(tt.input)

			var actual []interface{}
			SortedMapRange(m, tt.inputKind, func(key protoreflect.MapKey, value protoreflect.Value) bool {
				actual = append(actual, key.Interface())
				return true
			})

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("mismatch: want %+v, got %+v", tt.expected, actual)
			}
		})
	}
}

type mapT struct {
	protoreflect.Map
	inputMap reflect.Value
}

func newMapT(input interface{}) *mapT {
	return &mapT{
		inputMap: reflect.ValueOf(input),
	}
}

func (m mapT) Len() int {
	return m.inputMap.Len()
}

func (m mapT) Range(f func(protoreflect.MapKey, protoreflect.Value) bool) {
	for _, k := range m.inputMap.MapKeys() {
		f(protoreflect.MapKey(protoreflect.ValueOf(k.Interface())), protoreflect.Value{})
	}
}

func (m mapT) Get(key protoreflect.MapKey) protoreflect.Value {
	v := m.inputMap.MapIndex(reflect.ValueOf(key.Interface())).Interface()
	return protoreflect.ValueOf(v)
}
