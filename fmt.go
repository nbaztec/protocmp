package protocmp

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func quoteString(v interface{}) string {
	return strconv.Quote(fmt.Sprintf("%s", v))
}

func fmtMessage(m protoreflect.Message) string {
	f := &formatter{}
	f.printMessage(m)
	return f.String()
}

func fmtList(list protoreflect.List, fd protoreflect.FieldDescriptor) string {
	f := &formatter{}
	f.printList(list, fd)
	return f.String()
}

func fmtMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) string {
	f := &formatter{}
	f.printMap(mmap, fd)
	return f.String()
}

type formatter struct {
	str string
}

func (f *formatter) String() string {
	return strings.TrimSpace(f.str)
}

func (f *formatter) print(a ...interface{}) {
	for _, v := range a {
		f.str += fmt.Sprintf("%s", v)
	}
}

func (f *formatter) trim() {
	f.str = strings.TrimSuffix(f.str, " ")
}

func (f *formatter) trimAndPrint(v interface{}) {
	f.trim()
	f.str += fmt.Sprintf("%s", v)
}

func (f *formatter) printMessage(m protoreflect.Message) {
	f.print("<")
	defer f.trimAndPrint(">")
	
	messageDesc := m.Descriptor()
	fieldDescs := messageDesc.Fields()
	size := fieldDescs.Len()
	for i := 0; i < size; {
		fd := fieldDescs.Get(i)
		if od := fd.ContainingOneof(); od != nil {
			fd = m.WhichOneof(od)
			i += od.Fields().Len()
		} else {
			i++
		}

		if fd == nil || !m.Has(fd) {
			continue
		}

		name := fd.Name()
		f.print(name, ":")
		// Use type name for group field name.
		if fd.Kind() == protoreflect.GroupKind {
			name = fd.Message().Name()
		}
		val := m.Get(fd)
		f.printField(val, fd)

		f.print(" ")
	}
}

func (f *formatter) printField(val protoreflect.Value, fd protoreflect.FieldDescriptor) {
	switch {
	case fd.IsList():
		f.printList(val.List(), fd)
	case fd.IsMap():
		f.printMap(val.Map(), fd)
	default:
		f.printSingular(val, fd)
	}
}

func (f *formatter) printList(list protoreflect.List, fd protoreflect.FieldDescriptor) {
	f.print("[")
	defer f.trimAndPrint("]")

	size := list.Len()
	for i := 0; i < size; i++ {
		f.printSingular(list.Get(i), fd)
		if i != size-1 {
			f.print(" ")
		}
	}
}

func (f *formatter) printMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) {
	f.print("map[")
	defer f.trimAndPrint("]")

	SortedMapRange(mmap, fd.MapKey().Kind(), func(key protoreflect.MapKey, val protoreflect.Value) bool {
		f.printSingularKey(key.Value(), fd.MapKey())
		f.print(":")
		f.printSingular(val, fd.MapValue())
		f.print(" ")

		return true
	})
}

func (f *formatter) printSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) {
	f.printSingularWrapped(val, fd, true)
}

func (f *formatter) printSingularKey(val protoreflect.Value, fd protoreflect.FieldDescriptor) {
	f.printSingularWrapped(val, fd, false)
}

func (f *formatter) printSingularWrapped(val protoreflect.Value, fd protoreflect.FieldDescriptor, wrapString bool) {
	kind := fd.Kind()

	switch kind {
	case protoreflect.StringKind:
		if !wrapString {
			f.print(val.String())
			return
		}
		s := val.String()
		f.print("\"", s, "\"")

	case protoreflect.BoolKind, protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind,
		protoreflect.DoubleKind,
		protoreflect.BytesKind:
		f.print(val)

	case protoreflect.EnumKind:
		num := val.Enum()
		if desc := fd.Enum().Values().ByNumber(num); desc != nil {
			f.print(string(desc.Name()))
		} else {
			// Use numeric value if there is no enum description.
			f.print(int64(num))
		}

	case protoreflect.MessageKind, protoreflect.GroupKind:
		if !val.Message().IsValid() {
			f.print("<nil>")
			return
		}

		f.printMessage(val.Message())

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
}

