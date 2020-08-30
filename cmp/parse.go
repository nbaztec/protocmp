package cmp

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strings"
)

type fieldInfo struct {
	name      string
	fieldType reflect.Type
	kind      reflect.Kind
	value     interface{}
}

const (
	tagJSON          = "json"
	tagProtobuf      = "protobuf"
	tagProtobufOneof = "protobuf_oneof"
)

func parse(message proto.Message) []*fieldInfo {
	return parseRecursive(message, nil)
}

func parseRecursive(value interface{}, fields []*fieldInfo) []*fieldInfo {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	vt := reflect.TypeOf(value)
	if vt.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		vt = vt.Elem()
	}

	if isPrimitive(vt.Kind()) {
		return []*fieldInfo{
			{
				name:  "",
				kind:  vt.Kind(),
				value: v,
			},
		}
	}

	switch vt.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}

		var items []*fieldInfo
		for i := 0; i < v.Len(); i++ {
			slv := v.Index(i).Interface()
			slt := reflect.TypeOf(slv)
			if isPrimitive(slt.Kind()) {
				items = append(items, &fieldInfo{
					name:  fmt.Sprintf("[%d]", i),
					kind:  slt.Kind(),
					value: slv,
				})
			} else {
				fiItem := &fieldInfo{
					name: fmt.Sprintf("[%d]", i),
					kind: slt.Kind(),
				}

				if !reflect.ValueOf(slv).IsNil() {
					fiItem.kind = reflect.TypeOf(slv).Kind()
					fiItem.value = parseRecursive(slv, nil)
				}

				items = append(items, fiItem)
			}
		}
		return items
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		var items []*fieldInfo
		slSorted := Sort(v)
		for _, k := range slSorted.Key {
			slv := v.MapIndex(k).Interface()
			slt := reflect.TypeOf(slv)
			if isPrimitive(slt.Kind()) {
				items = append(items, &fieldInfo{
					name:  fmt.Sprintf("[%s]", k),
					kind:  slt.Kind(),
					value: slv,
				})
			} else {
				fiItem := &fieldInfo{
					name: fmt.Sprintf("[%s]", k),
					kind: slt.Kind(),
				}

				if !reflect.ValueOf(slv).IsNil() {
					fiItem.kind = reflect.TypeOf(slv).Kind()
					fiItem.value = parseRecursive(slv, nil)
				}

				items = append(items, fiItem)
			}
		}
		return items
	}

	//fmt.Println(v, v.Kind(), vt.Kind())
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)

		if !v.IsValid() {
			continue
		}

		// handle protobuf oneof fields
		_, ok := f.Tag.Lookup(tagProtobufOneof)
		if ok {
			fields = append(fields, parseRecursive(v.FieldByName(f.Name).Interface(), nil)...)
			continue
		}

		if !isPublicProtobufField(f) {
			continue
		}

		fv := v.FieldByName(f.Name).Interface()
		fvt := reflect.TypeOf(fv)

		fi := &fieldInfo{
			name:      f.Name,
			kind:      fvt.Kind(),
			fieldType: f.Type,
		}

		switch fvt.Kind() {
		case reflect.Ptr:
			fvv := reflect.ValueOf(fv)
			if !fvv.IsNil() {
				fi.value = parseRecursive(fvv.Elem().Interface(), nil)
			}

		case reflect.Slice:
			fiv := parseRecursive(fv, nil)
			if fiv != nil {
				fi.value = fiv
			}
		case reflect.Map:
			fiv := parseRecursive(fv, nil)
			if fiv != nil {
				fi.value = fiv
			}
		default:
			fi.value = fv
		}

		fields = append(fields, fi)
	}

	return fields
}

func isPublicProtobufField(f reflect.StructField) bool {
	jsonTag := f.Tag.Get(tagJSON)
	switch jsonTag {
	case "":
		tag := f.Tag.Get(tagProtobuf)
		return strings.Contains(tag, "json=")
	case "-":
		return false
	default:
		return true
	}
}

func isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return true
	default:
		return false
	}
}
