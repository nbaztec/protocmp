package cmp

import (
	"fmt"
	"reflect"
	"strings"
)

func fmtFunc(indent string, fields []*fieldInfo) string {
	s := ""
	for _, fi := range fields {
		switch fi.kind {
		case reflect.Ptr, reflect.Slice, reflect.Map:
			if fi.value == nil {
				s += fmt.Sprintf("%s%s: <nil>\n", indent, fi.name)
			} else {
				s += fmt.Sprintf("%s%s: \n", indent, fi.name)
				s += fmtFunc(indent+" ", fi.value.([]*fieldInfo))
			}
		default:
			s += fmt.Sprintf("%s%s: %+v\n", indent, fi.name, fi.value)
		}
	}

	return s
}

func fieldsToString(fields []*fieldInfo) string {
	return strings.TrimSuffix(fmtFunc("", fields), "\n")
}