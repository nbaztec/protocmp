package cmp

import (
	"fmt"
	"reflect"
	"strings"
)

func matchFields(expected, actual []*fieldInfo) *MatchError {
	return matchFieldsRecursive("", expected, actual, indexedFieldMapper())
}

func matchFieldsRecursive(key string, fieldsA, fieldsB []*fieldInfo, fm fieldMapper) *MatchError {
	intersect := map[string]int{}
	for i, fields := range [][]*fieldInfo{fieldsA, fieldsB} {
		for _, fi := range fields {
			if i == 0 {
				intersect[fi.name]++
			} else {
				intersect[fi.name]--
			}
		}
	}

	for name, v := range intersect {
		if v == 0 {
			continue
		}

		fieldKey := strings.TrimPrefix(key+"."+name, ".")
		if v == 1 {
			return newMatchError(fmt.Errorf("field invalid:\n+ %s\n- %s\n", fieldKey, fieldKey), fieldsA, fieldsB)
		} else {
			return newMatchError(fmt.Errorf("field invalid:\n- %s\n+ %s\n", fieldKey, fieldKey), fieldsA, fieldsB)
		}
	}

	for i := range fieldsA {
		fiA := fieldsA[i]
		fiB := fm(i, fiA.name, fieldsB)
		fieldKey := strings.TrimPrefix(key+"."+fiA.name, ".")

		if fiA.value == nil && fiB.value == nil {
			continue
		}

		if err := matchWithValueNil(fiA.value, fiB.value); err != nil {
			return newMatchError(fmt.Errorf("field value mismatch: %s\n%s\n", fieldKey, err), fieldsA, fieldsB)
		}

		switch fiA.kind {
		case reflect.Ptr:
			if err := matchFieldsRecursive(fieldKey, fiA.value.([]*fieldInfo), fiB.value.([]*fieldInfo), fm); err != nil {
				return err
			}

		case reflect.Map:
			fm = namedFieldMapper()
			fallthrough

		case reflect.Slice:
			fivA := fiA.value.([]*fieldInfo)
			fivB := fiB.value.([]*fieldInfo)
			if len(fivA) != len(fivB) {
				return newMatchError(fmt.Errorf("field length mismatch: %s\n+ %d\n- %d\n", fieldKey, len(fivA), len(fivB)), fivA, fivB)
			}

			if err := matchFieldsRecursive(fieldKey, fivA, fivB, fm); err != nil {
				return err
			}

		default:
			if fiA.kind != fiB.kind {
				return newMatchError(fmt.Errorf("field type mismatch: %s\n+ %+v\n- %+v\n", fieldKey, fiA.kind, fiB.kind), fieldsA, fieldsB)
			}

			if fiA.value != fiB.value {
				return newMatchError(fmt.Errorf("field value mismatch: %s\n+ %+v\n- %+v\n", fieldKey, fiA.value, fiB.value), fieldsA, fieldsB)
			}
		}
	}

	return nil
}


type fieldMapper func(index int, name string, fields []*fieldInfo) *fieldInfo

func indexedFieldMapper() fieldMapper {
	return func(index int, name string, fields []*fieldInfo) *fieldInfo {
		return fields[index]
	}
}

func namedFieldMapper() fieldMapper {
	return func(index int, name string, fields []*fieldInfo) *fieldInfo {
		for _, f := range fields {
			if f.name == name {
				return f
			}
		}

		return nil
	}
}

func matchWithValueNil(expected interface{}, actual interface{}) error {
	if expected != nil && actual != nil {
		return nil
	}

	if expected == nil {
		return fmt.Errorf("+ <nil>\n-\n%s", fieldsToString(actual.([]*fieldInfo)))
	}

	return fmt.Errorf("+\n%s\n- <nil>", fieldsToString(expected.([]*fieldInfo)))
}