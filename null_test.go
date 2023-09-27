package null_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/qawatake/go-sandbox/null"
)

func TestEqual(t *testing.T) {
	t.Run("non-Equaler", func(t *testing.T) {
		tests := []struct {
			name      string
			x1        null.T[int]
			x2        null.T[int]
			wantEqual bool
		}{
			{
				name:      "both are null",
				x1:        null.T[int]{},
				x2:        null.T[int]{},
				wantEqual: true,
			},
			{
				name:      "x1 is null",
				x1:        null.T[int]{},
				x2:        null.From(0),
				wantEqual: false,
			},
			{
				name:      "x2 is null",
				x1:        null.From(0),
				x2:        null.T[int]{},
				wantEqual: false,
			},
			{
				name:      "both are not null and equal",
				x1:        null.From(0),
				x2:        null.From(0),
				wantEqual: true,
			},
			{
				name:      "both are not null and not equal",
				x1:        null.From(0),
				x2:        null.From(1),
				wantEqual: false,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				assertEqual(t, tt.x1.Equal(tt.x2), tt.wantEqual)
			})
		}
	})

	t.Run("Equaler", func(t *testing.T) {
		tests := []struct {
			name      string
			x1        null.T[time.Time]
			x2        null.T[time.Time]
			wantEqual bool
		}{
			{
				name:      "both are null",
				x1:        null.T[time.Time]{},
				x2:        null.T[time.Time]{},
				wantEqual: true,
			},
			{
				name:      "x1 is null",
				x1:        null.T[time.Time]{},
				x2:        null.From(time.Time{}),
				wantEqual: false,
			},
			{
				name:      "x2 is null",
				x1:        null.From(time.Time{}),
				x2:        null.T[time.Time]{},
				wantEqual: false,
			},
			{
				name:      "both are not null and equal in the sense of ==",
				x1:        null.From(time.Time{}),
				x2:        null.From(time.Time{}),
				wantEqual: true,
			},
			{
				name: "both are not null and equal in the sense of Equal(time.Time) bool",
				// 2012-12-21T04:00:00Z
				x1: null.From(time.Date(2012, 12, 21, 4, 0, 0, 0, time.UTC)),
				// 2012-12-21T06:00:00+02:00
				x2:        null.From(time.Date(2012, 12, 21, 6, 0, 0, 0, time.FixedZone("", 2*60*60))),
				wantEqual: true,
			},
			{
				name:      "both are not null and not equal",
				x1:        null.From(time.Time{}),
				x2:        null.From(time.Date(2012, 12, 21, 21, 21, 21, 0, time.UTC)),
				wantEqual: false,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				assertEqual(t, tt.x1.Equal(tt.x2), tt.wantEqual)
			})
		}
	})
}

func TestScan(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		tests := []scanTestCase[bool]{
			{
				name:             format("true"),
				src:              true,
				wantValue:        true,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        false,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[bool]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("Float64", func(t *testing.T) {
		tests := []scanTestCase[float64]{
			{
				name:             format(1.2345),
				src:              1.2345,
				wantValue:        1.2345,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format("1.2345"),
				src:              "1.2345",
				wantValue:        1.2345,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[float64]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("Int64", func(t *testing.T) {
		tests := []scanTestCase[int]{
			{
				name:             format(int64(12345)),
				src:              int64(12345),
				wantValue:        12345,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[int]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		tests := []scanTestCase[string]{
			{
				name:             format("test"),
				src:              "test",
				wantValue:        "test",
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        "",
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[string]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("Time", func(t *testing.T) {
		const timeString1 = "2012-12-21T21:21:21Z"
		timeValue1, _ := time.Parse(time.RFC3339, timeString1)
		tests := []scanTestCase[time.Time]{
			{
				name:             format(timeValue1),
				src:              timeValue1,
				wantValue:        time.Date(2012, 12, 21, 21, 21, 21, 0, time.UTC),
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:      format(int64(42)),
				src:       int64(42),
				wantValue: time.Time{},
				// database/sql Scan set Valid to true in this case.
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[time.Time]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("CustomScanner", func(t *testing.T) {

		tests := []scanTestCase[customScanner]{
			{
				name:             format(true),
				src:              true,
				wantValue:        customScanner{"true"},
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(int64(12345)),
				src:              int64(12345),
				wantValue:        customScanner{"12345"},
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format(nil),
				src:              nil,
				wantValue:        customScanner{""},
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[customScanner]
				err := nullable.Scan(tt.src)
				tt.requireErrorFunc(t, err)
				assertEqualStruct(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[bool]{
			{
				name:             format([]byte("true")),
				data:             []byte("true"),
				wantValue:        true,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`{"Bool":true,"Valid":true}`)),
				data:             []byte(`{"Bool":true,"Valid":true}`),
				wantValue:        false,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        false,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`12345`)),
				data:             []byte(`12345`),
				wantValue:        false,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        false,
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[bool]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("Float64", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[float64]{
			{
				name:             format([]byte(`1.2345`)),
				data:             []byte(`1.2345`),
				wantValue:        1.2345,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name: format([]byte(`"1.2345"`)),
				data: []byte([]byte(`"1.2345"`)),
				// difference from guregu/null
				// wantValue:        1.2345,
				// wantIsNull:       false,
				// requireErrorFunc: requireNoError,
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`{"Float64":1.2345,"Valid":true}`)),
				data:             []byte(`{"Float64":1.2345,"Valid":true}`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`""`)),
				data:             []byte(`""`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`true`)),
				data:             []byte(`true`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[float64]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})
	t.Run("Int64", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[int64]{
			{
				name:             format([]byte(`12345`)),
				data:             []byte(`12345`),
				wantValue:        12345,
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name: format([]byte(`"12345"`)),
				data: []byte(`"12345"`),
				// difference from guregu/null
				// wantValue:        12345,
				// wantIsNull:       false,
				// requireErrorFunc: requireNoError,
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`{"Int64":12345,"Valid":true}`)),
				data:             []byte(`{"Int64":12345,"Valid":true}`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`""`)),
				data:             []byte(`""`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`true`)),
				data:             []byte(`true`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        0,
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[int64]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[string]{
			{
				name:             format([]byte(`"test"`)),
				data:             []byte(`"test"`),
				wantValue:        "test",
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`{"String":"test","Valid":true}`)),
				data:             []byte(`{"String":"test","Valid":true}`),
				wantValue:        "",
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`""`)),
				data:             []byte(`""`),
				wantValue:        "",
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        "",
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`true`)),
				data:             []byte(`true`),
				wantValue:        "",
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        "",
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[string]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})

	t.Run("Time", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[time.Time]{
			{
				name:             format([]byte(`"2012-12-21T21:21:21Z"`)),
				data:             []byte(`"2012-12-21T21:21:21Z"`),
				wantValue:        time.Date(2012, 12, 21, 21, 21, 21, 0, time.UTC),
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`)),
				data:             []byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)),
				data:             []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
			{
				name:             "bad object" + format([]byte(`{"hello": "world"}`)),
				data:             []byte(`{"hello": "world"}`),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`12345`)),
				data:             []byte(`12345`),
				wantValue:        time.Time{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[time.Time]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqual(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})
	t.Run("CustomJSONUnmarshaler", func(t *testing.T) {
		tests := []unmarshalJSONTestCase[customJSONUnmarshaler]{
			{
				name:             format([]byte(`12345`)),
				data:             []byte(`12345`),
				wantValue:        customJSONUnmarshaler{12345},
				wantIsNull:       false,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`"12345"`)),
				data:             []byte(`"12345"`),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`{"i":12345,"Valid":true}`)),
				data:             []byte(`{"i":12345,"Valid":true}`),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte(`""`)),
				data:             []byte(`""`),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             format([]byte("null")),
				data:             []byte("null"),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireNoError,
			},
			{
				name:             format([]byte(`true`)),
				data:             []byte(`true`),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireError,
			},
			{
				name:             "invalid json",
				data:             []byte(`:)`),
				wantValue:        customJSONUnmarshaler{},
				wantIsNull:       true,
				requireErrorFunc: requireJSONSyntaxError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				var nullable null.T[customJSONUnmarshaler]
				err := json.Unmarshal(tt.data, &nullable)
				tt.requireErrorFunc(t, err)
				assertEqualStruct(t, nullable.ValueOrZero(), tt.wantValue)
				assertEqual(t, nullable.IsNull(), tt.wantIsNull)
			})
		}
	})
}

func TestUnmarshalJSON_AsField(t *testing.T) {
	type Object struct {
		NullableInt null.T[int]
	}

	tests := []struct {
		name       string
		data       []byte
		wantIsNull bool
		wantValue  int
		requireErrorFunc
	}{
		{
			name:             "explicit null",
			data:             []byte(`{"NullableInt":null}`),
			wantIsNull:       true,
			wantValue:        0,
			requireErrorFunc: requireNoError,
		},
		{
			name:             "omitted",
			data:             []byte(`{}`),
			wantIsNull:       true,
			wantValue:        0,
			requireErrorFunc: requireNoError,
		},
		{
			name:             "zero",
			data:             []byte(`{"NullableInt":0}`),
			wantIsNull:       false,
			wantValue:        0,
			requireErrorFunc: requireNoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var obj Object
			err := json.Unmarshal(tt.data, &obj)
			tt.requireErrorFunc(t, err)
			assertEqual(t, obj.NullableInt.IsNull(), tt.wantIsNull)
			assertEqual(t, obj.NullableInt.ValueOrZero(), tt.wantValue)
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Run("Float", func(t *testing.T) {
		tests := []marshalJSONTestCase[float64]{
			{
				name:             format(1.2345),
				src:              1.2345,
				wantData:         []byte(`1.2345`),
				requireErrorFunc: requireNoError,
			},
			{
				name: format(-9.40623162845385e-07),
				src:  -9.40623162845385e-07,
				// difference from guregu/null
				// guregu/null.Float.MarshalJSON() returns []byte(`-0.000000940623162845385`)
				// "-0.000000940623162845385",
				wantData:         []byte(`-9.40623162845385e-7`),
				requireErrorFunc: requireNoError,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				f := null.From(tt.src)
				data, err := f.MarshalJSON()
				tt.requireErrorFunc(t, err)
				assertEqual(t, string(data), string(tt.wantData))
			})
		}
	})
}

func TestMarshalJSON_AsField(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		type Object struct {
			NullableInt null.T[int]
		}
		obj := Object{
			NullableInt: null.From(1),
		}
		b, err := json.Marshal(obj)
		requireNoError(t, err)
		assertEqual(t, string(b), `{"NullableInt":1}`)
	})

	t.Run("null without omitempty tag", func(t *testing.T) {
		type Object struct {
			NullableInt null.T[int]
		}
		obj := Object{
			NullableInt: null.T[int]{},
		}
		b, err := json.Marshal(obj)
		requireNoError(t, err)
		assertEqual(t, string(b), `{"NullableInt":null}`)
	})

	// zero values of structs with omitempty tags are not omitted.
	// https://github.com/golang/go/issues/11939
	t.Run("null with omitempty tag", func(t *testing.T) {
		type Object struct {
			NullableInt null.T[int] `json:",omitempty"`
		}
		obj := Object{
			NullableInt: null.T[int]{},
		}
		b, err := json.Marshal(obj)
		requireNoError(t, err)
		assertEqual(t, string(b), `{"NullableInt":null}`)
	})
}

func TestPtr(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		i := null.From[int](100)
		p := i.Ptr()
		assertEqual(t, *p, 100)
	})

	t.Run("null", func(t *testing.T) {
		i := null.T[int]{}
		p := i.Ptr()
		assertEqual(t, p, nil)
	})
}

func TestFromPtr(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		var i uint = 100
		n := null.FromPtr(&i)
		assertEqual(t, n.IsNull(), false)
		assertEqual(t, n.ValueOrZero(), uint(100))
	})

	t.Run("null", func(t *testing.T) {
		i := null.FromPtr[uint](nil)
		assertEqual(t, i.IsNull(), true)
		assertEqual(t, i.ValueOrZero(), 0)
	})
}

func Test_SharedValues(t *testing.T) {
	t.Run("Scan (Immutable)", func(t *testing.T) {
		i := null.From[int](100)
		j := i
		err := j.Scan(200)
		requireNoError(t, err)

		assertEqual(t, i.ValueOrZero(), 100)
		assertEqual(t, j.ValueOrZero(), 200)
	})

	t.Run("Scan (Pointer)", func(t *testing.T) {
		p := toptr(100)
		i := null.From(p)
		j := i
		err := j.Scan(200)
		requireNoError(t, err)

		assertEqual(t, i.ValueOrZero(), toptr(100))
		assertEqual(t, j.ValueOrZero(), toptr(200))
	})

	t.Run("Scan Failure", func(t *testing.T) {
		i := null.From(neverScanner{100})
		j := i
		err := j.Scan("200")
		requireError(t, err)

		assertEqualStruct(t, i.ValueOrZero(), neverScanner{100})
		assertEqualStruct(t, j.ValueOrZero(), neverScanner{})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		i := null.From[int](100)
		j := i
		err := j.UnmarshalJSON([]byte("200"))
		requireNoError(t, err)

		assertEqual(t, i.ValueOrZero(), 100)
		assertEqual(t, j.ValueOrZero(), 200)
	})

	t.Run("UnmarshalJSON (Pointer)", func(t *testing.T) {
		p := toptr(100)
		i := null.From(p)
		j := i
		err := j.UnmarshalJSON([]byte("200"))
		requireNoError(t, err)

		assertEqual(t, i.ValueOrZero(), toptr(100))
		assertEqual(t, j.ValueOrZero(), toptr(200))
	})

	t.Run("UnmarshalJSON Failure", func(t *testing.T) {
		i := null.From(neverJSONUnmarshaler{100})
		j := i
		err := j.UnmarshalJSON([]byte("200"))
		requireError(t, err)

		assertEqualStruct(t, i.ValueOrZero(), neverJSONUnmarshaler{100})
		assertEqualStruct(t, j.ValueOrZero(), neverJSONUnmarshaler{})
	})

	t.Run("Ptr", func(t *testing.T) {
		i := null.From[int](100)
		p0 := i.Ptr()
		p1 := i.Ptr()
		p2 := i.Ptr()
		*p1 = 200
		*p2 = 400

		assertEqual(t, *p0, 100)
		assertEqual(t, *p1, 200)
		assertEqual(t, *p2, 400)
	})
}

type scanTestCase[C comparable] struct {
	name       string
	src        any
	wantValue  C
	wantIsNull bool
	requireErrorFunc
}

type customScanner struct {
	s string
}

func (s *customScanner) Scan(src any) error {
	s.s = fmt.Sprint(src)
	return nil
}

type neverScanner struct {
	i int
}

func (s *neverScanner) Scan(src any) error {
	return errors.New("never")
}

type unmarshalJSONTestCase[C comparable] struct {
	name       string
	data       []byte
	wantValue  C
	wantIsNull bool
	requireErrorFunc
}

type customJSONUnmarshaler struct {
	i int
}

func (x *customJSONUnmarshaler) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	x.i = i
	return nil
}

type neverJSONUnmarshaler struct {
	i int
}

func (x *neverJSONUnmarshaler) UnmarshalJSON(data []byte) error {
	return errors.New("never")
}

type marshalJSONTestCase[T any] struct {
	name     string
	src      T
	wantData []byte
	requireErrorFunc
}

type requireErrorFunc func(t *testing.T, err error)

var _ requireErrorFunc = requireError

func requireError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("want error, but got nil")
	}
}

func requireJSONSyntaxError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("want error, but got nil")
	}
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

var _ requireErrorFunc = requireNoError

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
}

func assertEqual[T comparable](t *testing.T, x T, y T, opts ...cmp.Option) bool {
	t.Helper()
	if diff := cmp.Diff(x, y, opts...); diff != "" {
		t.Errorf(diff)
		return false
	}
	return true
}

func assertEqualStruct[T comparable](t *testing.T, x, y T) bool {
	t.Helper()
	var z T
	return assertEqual[T](t, x, y, cmp.AllowUnexported(z))
}

func format(v any) string {
	switch v := v.(type) {
	case []byte:
		return fmt.Sprintf("%T:%q", v, v)
	case string:
		return fmt.Sprintf("%T:%q", v, v)
	case fmt.Stringer:
		return fmt.Sprintf("%T:%q", v, v)
	default:
		return fmt.Sprintf("%T:%v", v, v)
	}
}

func toptr[T any](x T) *T {
	return &x
}
