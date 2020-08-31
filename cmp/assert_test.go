package cmp

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

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

	actual := &sample.Outer{
		StrVal:    "foo",
		IntVal:    1,
		BoolVal:   true,
		DoubleVal: 1.1,
		BytesVal:  []byte{0x01, 0x02},
		RepeatedType: []*sample.Outer_Inner{
			{Id: "1"},
			{Id: "2"},
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

	AssertEqual(t, expected, actual)
}
