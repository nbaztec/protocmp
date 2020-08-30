# protocmp
Compare protobuf messages and print diff

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
     TestCmpAssertEqual: assert.go:18: field value mismatch IntVal
         + 1
         - 10
         
         ++ expected
         StrVal: foo
         IntVal: 1
         BoolVal: true
         DoubleVal: 1.1
         BytesVal: 
          [0]: 1
          [1]: 2
         RepeatedType: 
          [0]: 
           Id: 1
          [1]: 
           Id: 2
          [2]: <nil>
         MapType: 
          [A]: 
           Id: AA
          [B]: 
           Id: BB
          [C]: <nil>
         EnumType: NOT_OK
         OneofInt: 1
         LastUpdated: 
          Seconds: 1598794689
          Nanos: 562367932
         LastUpdatedDuration: 
          Seconds: 1
          Nanos: 0
         Details: 
          TypeUrl: mytype/v1
          Value: 
           [0]: 5
         RepeatedTypeSimple: 
          [0]: 9
          [1]: 10
          [2]: 11
         
         -- actual
         StrVal: foo
         IntVal: 10
         BoolVal: true
         DoubleVal: 1.1
         BytesVal: 
          [0]: 1
          [1]: 2
         RepeatedType: 
          [0]: 
           Id: 1
          [1]: 
           Id: 3
          [2]: <nil>
         MapType: 
          [C]: <nil>
          [A]: 
           Id: AA
          [B]: 
           Id: BB
         EnumType: NOT_OK
         OneofInt: 1
         LastUpdated: 
          Seconds: 1598794689
          Nanos: 562367932
         LastUpdatedDuration: 
          Seconds: 1
          Nanos: 0
         Details: 
          TypeUrl: mytype/v1
          Value: 
           [0]: 5
         RepeatedTypeSimple: 
          [0]: 9
          [1]: 10
          [2]: 11
 --- FAIL: TestCmpAssertEqual (0.00s)
 FAIL
```