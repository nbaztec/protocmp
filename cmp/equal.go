package cmp

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Equal(x, y proto.Message) error {
	if err := equal(proto.MessageV2(x), proto.MessageV2(y)); err != nil {
		return err
	}

	return nil
}

func equal(x, y protoreflect.ProtoMessage) *matchErr {
	if x == nil || y == nil {
		if x == nil && y == nil {
			return nil
		}

		return newMatchError("value mismatch").Values(x, y)
	}

	mx := x.ProtoReflect()
	my := y.ProtoReflect()
	if mx.IsValid() != my.IsValid() {
		return newMatchError("invalid value").Values(x, y)
	}
	if err := equalMessage(mx, my); err != nil {
		return err
	}

	return nil
}

func (m *matchErr) Error() string {
	return fmt.Sprintf("%s: %s\n+ %v\n- %v", strings.Join(m.fieldKeys, "."), m.message, m.expected, m.actual)
}


func fmtError(v protoreflect.Value, fd protoreflect.FieldDescriptor) *matchErr {
	pfmt := newFormatter()
	switch {
	case fd.IsList():
		pfmt.printList("", v.List(), fd)
		return newMatchError("value mismatch").Field(fd.Name()).Values(pfmt.String(), nil)
	case fd.IsMap():
		pfmt.printMap("", v.Map(), fd)
		return newMatchError("value mismatch").Field(fd.Name()).Values(pfmt.String(), nil)
	default:
		return newMatchError("value mismatch").Field(fd.Name()).Values(v.Interface(), nil)
	}
}

// equalMessage compares two messages.
func equalMessage(mx, my protoreflect.Message) *matchErr {
	if mx.Descriptor() != my.Descriptor() {
		return newMatchError("descriptors don't match")
	}

	if mx.IsValid() && !my.IsValid() {
		pfmt := newFormatter()
		pfmt.printMessage(mx)
		return newMatchError("value mismatch").Values(pfmt.String(), nil)
	}

	if !mx.IsValid() && my.IsValid() {
		return newMatchError("value mismatch").Values(nil, my)
	}


	//if !mx.IsValid() && my.IsValid() {
	//	return newMatchError("value mismatch").Values(nil, fmt.Sprintf("%v (*%s)", my.Interface(), my.Descriptor().Name()))
	//}
	//
	//if mx.IsValid() && !my.IsValid() {
	//	return newMatchError("value mismatch").Values(fmt.Sprintf("%v (*%s)", mx.Interface(), mx.Descriptor().Name()), nil)
	//}

	nx := 0
	var equalErr *matchErr
	mx.Range(func(fd protoreflect.FieldDescriptor, vx protoreflect.Value) bool {
		nx++
		vy := my.Get(fd)

		switch  {
		case fd.IsList() || fd.IsMap():
			fallthrough
		case fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind:
			if !my.Has(fd) {
				equalErr = fmtError(vx, fd).ValueActual(nil)
				return false
			}
		default:
			if !my.Has(fd) {
				equalErr = fmtError(vx, fd).ValueActual(vy.Interface())
				return false
			}
		}
		//switch fd.Name() {
		//case "repeated_type":
		//	pfmt := newFormatter()
		//	pfmt.printList(string(fd.Name()), vx.List(), fd)
		//	fmt.Println(pfmt)
		//case "repeated_type_simple":
		//	pfmt := newFormatter()
		//	pfmt.printList(string(fd.Name()), vx.List(), fd)
		//	fmt.Println(pfmt)
		//case "map_type":
		//	pfmt := newFormatter()
		//	pfmt.printMap(string(fd.Name()), vx.Map(), fd)
		//	fmt.Println(pfmt)
		//}



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
	my.Range(func(fd protoreflect.FieldDescriptor, vy protoreflect.Value) bool {
		ny++
		if !mx.Has(fd) {
			equalErr = fmtError(vy, fd).Values(nil, vy)
			return false
		}

		return true
	})

	if equalErr != nil {
		return equalErr
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
		return newMatchError("length mismatch").Values(x.Len(), y.Len())
	}
	var equalErr *matchErr
	x.Range(func(k protoreflect.MapKey, vx protoreflect.Value) bool {
		vy := y.Get(k)
		if !y.Has(k) {
			equalErr = newMatchError("missing key").Field(protoreflect.Name(fmt.Sprintf("[%s]", k.String())))
			return false
		}
		if err := equalValue(fd.MapValue(), vx, vy); err != nil {
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
		return newMatchError("length mismatch").Values(x.Len(), y.Len())
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
		//xm := x.Message()
		//ym := y.Message()
		//
		//if xm.IsValid() && !ym.IsValid() {
		//	pfmt := newFormatter()
		//	pfmt.printMessage(x.Message())
		//	return newMatchError("value mismatch").Values(pfmt.String(), nil)
		//}
		//
		//if !xm.IsValid() && ym.IsValid() {
		//	return newMatchError("value mismatch").Values(nil, y)
		//}

		if err := equalMessage(x.Message(), y.Message()); err != nil {
			return err
		}
		return nil
	case fd.Kind() == protoreflect.BytesKind:
		if !bytes.Equal(x.Bytes(), y.Bytes()) {
			return newMatchError("value mismatch").Values(x.Bytes(), y.Bytes())
		}
	case fd.Kind() == protoreflect.FloatKind, fd.Kind() == protoreflect.DoubleKind:
		fx := x.Float()
		fy := y.Float()
		if math.IsNaN(fx) || math.IsNaN(fy) {
			if !math.IsNaN(fx) && math.IsNaN(fy) {
				return newMatchError("value mismatch").Values(fx, fy)
			}
			if math.IsNaN(fx) && !math.IsNaN(fy) {
				return newMatchError("value mismatch").Values(fx, fy)
			}

			return nil
		}
		if fx != fy {
			return newMatchError("value mismatch").Values(fx, fy)
		}

		return nil
	default:
		if x.Interface() != y.Interface() {
			return newMatchError("value mismatch").Values(x.Interface(), y.Interface())
		}

		return nil
	}

	return nil
}

// equalUnknown compares unknown fields by direct comparison on the raw bytes
// of each individual field number.
func equalUnknown(x, y protoreflect.RawFields) *matchErr {
	if len(x) != len(y) {
		return newMatchError("length mismatch").Values(len(x), len(y))
	}
	if !bytes.Equal(x, y) {
		return newMatchError("value mismatch").Values(x, y)
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
		return newMatchError("value mismatch").Values(mx, my)
	}

	return nil
}
