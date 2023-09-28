package null_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/qawatake/null"
)

func ExampleT() {
	var i1, i2 null.T[int]
	fmt.Printf("i1 == i2: %v\n", i1 == i2)

	var a1, a2 null.T[[3]bool]
	fmt.Printf("a1 == a2: %v\n", a1 == a2)
	var s1, s2 null.T[struct{ a int }]
	fmt.Printf("s1 == s2: %v\n", s1 == s2)

	// panic: runtime error
	// x1 := null.From[any](map[string]int{})
	// x2 := null.From[any](map[string]int{})
	// fmt.Printf("x1 == x2: %v\n", x1 == x2)

	// Output:
	// i1 == i2: true
	// a1 == a2: true
	// s1 == s2: true
}

// https://go.dev/blog/comparable#consequences-and-remedies
func isComparable[_ comparable]() {}

// null.T[int] is strictly comparable.
func _[P null.T[int]]() {
	_ = isComparable[P]
}

// null.T[[2]int] is strictly comparable.
func _[P null.T[[2]int]]() {
	_ = isComparable[P]
}

// null.T[any] is not strictly comparable.
func _[P null.T[any]]() {
	// _ = isComparable[P]
}

func ExampleT_Scan() {
	var age null.T[time.Duration]
	age.Scan(int64(24 * 1000 * time.Hour))

	fmt.Printf("valid: %v\n", !age.IsNull())
	fmt.Printf("type: %T\n", age.ValueOrZero())
	fmt.Printf("value: %v\n", age.ValueOrZero())
	// Output:
	// valid: true
	// type: time.Duration
	// value: 24000h0m0s
}

func ExampleT_UnmarshalJSON() {
	type ComparableStruct struct {
		Bool   bool
		Int    int64
		Float  float64
		String string
	}
	type ComparableArray [4]int
	type Object struct {
		Duration null.T[time.Duration]
		Struct   null.T[ComparableStruct]
		Array    null.T[ComparableArray]
	}

	var obj1 Object
	data1 := []byte(`{
		"Duration": 1000,
		"Struct": {
			"Bool": true,
			"Int": 123,
			"Float": 1.23,
			"String": "abc"
		},
		"Array": [1,2,3]
	}`)
	json.Unmarshal(data1, &obj1)
	fmt.Println("[obj1]")
	fmt.Printf("duration:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj1.Duration.IsNull(), obj1.Duration.ValueOrZero())
	fmt.Printf("struct:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj1.Struct.IsNull(), obj1.Struct.ValueOrZero())
	fmt.Printf("array:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj1.Array.IsNull(), obj1.Array.ValueOrZero())

	var obj2 Object
	data2 := []byte(`{}`)
	json.Unmarshal(data2, &obj2)
	fmt.Println("[obj2]")
	fmt.Printf("duration:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj2.Duration.IsNull(), obj2.Duration.ValueOrZero())
	fmt.Printf("struct:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj2.Struct.IsNull(), obj2.Struct.ValueOrZero())
	fmt.Printf("array:\n  valid: %+[1]v\n  type: %[2]T\n  value: %+[2]v\n", !obj2.Array.IsNull(), obj2.Array.ValueOrZero())
	// Output:
	// [obj1]
	// duration:
	//   valid: true
	//   type: time.Duration
	//   value: 1µs
	// struct:
	//   valid: true
	//   type: null_test.ComparableStruct
	//   value: {Bool:true Int:123 Float:1.23 String:abc}
	// array:
	//   valid: true
	//   type: null_test.ComparableArray
	//   value: [1 2 3 0]
	// [obj2]
	// duration:
	//   valid: false
	//   type: time.Duration
	//   value: 0s
	// struct:
	//   valid: false
	//   type: null_test.ComparableStruct
	//   value: {Bool:false Int:0 Float:0 String:}
	// array:
	//   valid: false
	//   type: null_test.ComparableArray
	//   value: [0 0 0 0]
}

func ExampleT_Equal() {
	// 2012-12-21T04:00:00Z
	x1 := null.From(time.Date(2012, 12, 21, 4, 0, 0, 0, time.UTC))
	// 2012-12-21T06:00:00+02:00
	x2 := null.From(time.Date(2012, 12, 21, 6, 0, 0, 0, time.FixedZone("", 2*60*60)))
	var x3 null.T[time.Time]

	fmt.Printf("x1 == x2: %v\n", x1 == x2)
	fmt.Printf("x1.Equal(x2): %v\n", x1.Equal(x2))
	fmt.Printf("x1.Equal(x3): %v\n", x3.Equal(x1))
	// Output:
	// x1 == x2: false
	// x1.Equal(x2): true
	// x1.Equal(x3): false
}

func ExampleT_Ptr() {
	i := null.From[int](123)
	p1 := i.Ptr()
	p2 := i.Ptr()
	fmt.Printf("p1 != p2: %v\n", p1 != p2)
	// Output:
	// p1 != p2: true
}
