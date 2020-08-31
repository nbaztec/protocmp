package cmplegacy

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/nbaztec/protocmp/protos/sample"
)

func TestBar(t *testing.T) {
	now, _ := ptypes.TimestampProto(time.Date(2020, time.August, 30, 19, 05, 00, 00, time.UTC))
	expected := &sample.Outer{
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

	actual := &sample.Outer{
		StrVal:    "foo",
		IntVal:    1,
		BoolVal:   true,
		DoubleVal: 1.1,
		BytesVal:  []byte{0x01, 0x02},
		RepeatedType: []*sample.Outer_Inner{
			{Id: "1"},
			{Id: "2"},
			{Id: "2"},
			//nil,
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


	fmt.Println(EqualV2(proto.MessageV2(expected), proto.MessageV2(actual)))
}


func EqualV2(x, y protoreflect.ProtoMessage) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	mx := x.ProtoReflect()
	my := y.ProtoReflect()
	if mx.IsValid() != my.IsValid() {
		return false
	}
	if err:= equalMessage(mx, my); err !=nil {
		fmt.Println(err)
		return false
	}

	return true
}

type matchErr struct {
	fieldKeys []string
	message string
	expected interface{}
	actual interface{}
}

func newMatchV2Error(message string) *matchErr {
	return &matchErr{ message: message}
}

func (m *matchErr) Field(k protoreflect.Name) *matchErr {
	m.fieldKeys = append([]string{string(k)}, m.fieldKeys...)
	return m
}

func (m *matchErr) Values(expected, actual interface{}) *matchErr {
	m.expected = expected
	m.actual = actual
	return m
}

func (m *matchErr) Error() string {
	return fmt.Sprintf("%s: %s\n+ %v\n- %v", strings.Join(m.fieldKeys, "."), m.message, m.expected, m.actual)
}

// equalMessage compares two messages.
func equalMessage(mx, my protoreflect.Message) *matchErr {
	if mx.Descriptor() != my.Descriptor() {
		return newMatchV2Error("descriptors don't match")
	}

	fmt.Println(mx.Descriptor().Name(), mx.Interface())
	if !mx.IsValid() && my.IsValid() {
		return newMatchV2Error("value mismatch").Values(nil, fmt.Sprintf("%v (*%s)", my.Interface(), my.Descriptor().Name()))
	}
	if mx.IsValid() && !my.IsValid() {
		return newMatchV2Error("value mismatch").Values(fmt.Sprintf("%v (*%s)", mx.Interface(), mx.Descriptor().Name()), nil)
	}

	intersectFields := map[protoreflect.Name]int{}

	nx := 0
	var equalErr *matchErr
	mx.Range(func(fd protoreflect.FieldDescriptor, vx protoreflect.Value) bool {
		nx++
		intersectFields[fd.Name()]++
		vy := my.Get(fd)
		if !my.Has(fd) {
			equalErr = newMatchV2Error("missing field").Field(fd.Name())
			return false
		}
		if err := equalField(fd, vx, vy); err != nil {
			equalErr = err
			return false
		}

		return true
	})

	if equalErr != nil {
		return equalErr
	}
	ny := 0
	my.Range(func(fd protoreflect.FieldDescriptor, vx protoreflect.Value) bool {
		ny++
		intersectFields[fd.Name()]--
		return true
	})
	for name, v := range intersectFields {
		if v == 0 {
			continue
		}

		return newMatchV2Error("extra field").Field(name)
	}

	return equalUnknown(mx.GetUnknown(), my.GetUnknown())
}

// equalField compares two fields.
func equalField(fd protoreflect.FieldDescriptor, x, y protoreflect.Value) *matchErr {
	switch {
	case fd.IsList():
		if err := equalList(fd, x.List(), y.List()); err != nil {
			return err.Field(fd.Name())
		}
	case fd.IsMap():
		if err := equalMap(fd, x.Map(), y.Map()); err != nil {
			return err.Field(fd.Name())
		}
	default:
		if err := equalValue(fd, x, y); err != nil {
			return err.Field(fd.Name())
		}
	}

	return nil
}

// equalMap compares two maps.
func equalMap(fd protoreflect.FieldDescriptor, x, y protoreflect.Map) *matchErr {
	if x.Len() != y.Len() {
		return newMatchV2Error("length mismatch").Values(x.Len(), y.Len())
	}
	var equalErr *matchErr
	x.Range(func(k protoreflect.MapKey, vx protoreflect.Value) bool {
		vy := y.Get(k)
		if !y.Has(k) {
			equalErr = newMatchV2Error("missing key").Field(protoreflect.Name(fmt.Sprintf("[%s]", k.String())))
			return false
		}
		if err := equalValue(fd.MapValue(), vx, vy); err !=nil {
			equalErr = err.Field(protoreflect.Name(fmt.Sprintf("[%s]", k.String())))
			return false
		}
		return true
	})
	return equalErr
}

// equalList compares two lists.
func equalList(fd protoreflect.FieldDescriptor, x, y protoreflect.List) *matchErr {
	if x.Len() != y.Len() {
		return newMatchV2Error("length mismatch").Values(x.Len(), y.Len())
	}
	for i := x.Len() - 1; i >= 0; i-- {
		if err := equalValue(fd, x.Get(i), y.Get(i)); err != nil {
			return err.Field(protoreflect.Name(fmt.Sprintf("[%d]", i)))
		}
	}
	return nil
}

// equalValue compares two singular values.
func equalValue(fd protoreflect.FieldDescriptor, x, y protoreflect.Value) *matchErr {
	switch {
	case fd.Message() != nil:
		if err := equalMessage(x.Message(), y.Message()); err != nil {
			return err
		}
		return nil
	case fd.Kind() == protoreflect.BytesKind:
		if !bytes.Equal(x.Bytes(), y.Bytes()) {
			return newMatchV2Error("value mismatch").Values(x.Bytes(), y.Bytes())
		}
	case fd.Kind() == protoreflect.FloatKind, fd.Kind() == protoreflect.DoubleKind:
		fx := x.Float()
		fy := y.Float()
		if math.IsNaN(fx) || math.IsNaN(fy) {
			if !math.IsNaN(fx) && math.IsNaN(fy) {
				return newMatchV2Error("value mismatch").Values(fx, fy)
			}
			if math.IsNaN(fx) && !math.IsNaN(fy) {
				return newMatchV2Error("value mismatch").Values(fx, fy)
			}

			return nil
		}
		if fx != fy {
			return newMatchV2Error("value mismatch").Values(fx, fy)
		}

		return nil
	default:
		if x.Interface() != y.Interface() {
			return newMatchV2Error("value mismatch").Values(x.Interface(), y.Interface())
		}

		return nil
	}

	return nil
}

// equalUnknown compares unknown fields by direct comparison on the raw bytes
// of each individual field number.
func equalUnknown(x, y protoreflect.RawFields) *matchErr {
	if len(x) != len(y) {
		return newMatchV2Error("length mismatch").Values(len(x), len(y))
	}
	if !bytes.Equal(x, y) {
		return newMatchV2Error("value mismatch").Values(x, y)
	}

	mx := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	my := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	for len(x) > 0 {
		fnum, _, n := protowire.ConsumeField(x)
		mx[fnum] = append(mx[fnum], x[:n]...)
		x = x[n:]
	}
	for len(y) > 0 {
		fnum, _, n := protowire.ConsumeField(y)
		my[fnum] = append(my[fnum], y[:n]...)
		y = y[n:]
	}
	if !reflect.DeepEqual(mx, my) {
		return newMatchV2Error("value mismatch").Values(mx, my)
	}

	return nil
}
