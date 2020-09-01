# protocmp
Compare protobuf messages and compute diff

### Usage
```go
func TestFoo(t *testing.T) {
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
    		IntVal:    10,
    		BoolVal:   true,
    		DoubleVal: 1.1,
    		BytesVal:  []byte{0x01, 0x02},
    		RepeatedType: []*sample.Outer_Inner{
    			{Id: "1"},
    			{Id: "3"},
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
    
    	cmp.AssertEqual(t, a, b)
}
```
```
/=== RUN   TestCmpAssertEqual
     TestCmpAssertEqual: assert.go:18: int_val: value mismatch
         + 1
         - 10
 --- FAIL: TestCmpAssertEqual (0.00s)
 FAIL
```