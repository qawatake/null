package null_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/qawatake/go-sandbox/null"
	gnull "gopkg.in/guregu/null.v4"
)

func FuzzInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, in int64) {
		n0 := gnull.IntFrom(in)
		n1 := null.From[int64](in)

		if diff := cmp.Diff(n0.IsZero(), n1.IsNull()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.ValueOrZero(), n1.ValueOrZero()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.Ptr(), n1.Ptr()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		b0, err1 := n0.MarshalJSON()
		b1, err2 := n1.MarshalJSON()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}
		if diff := cmp.Diff(string(b0), string(b1)); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		v0, err1 := n0.Value()
		v1, err2 := n1.Value()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}
		if diff := cmp.Diff(v0, v1); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}
	})
}

func FuzzFloat(f *testing.F) {
	f.Fuzz(func(t *testing.T, in float64) {
		n0 := gnull.FloatFrom(in)
		n1 := null.From[float64](in)

		if diff := cmp.Diff(n0.IsZero(), n1.IsNull()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.ValueOrZero(), n1.ValueOrZero()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.Ptr(), n1.Ptr()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if !isFloatMarshalJSONIgnoreCase(t, in) {
			b0, err1 := n0.MarshalJSON()
			b1, err2 := n1.MarshalJSON()
			if !errors.Is(err2, err1) {
				t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
			}
			if diff := cmp.Diff(string(b0), string(b1)); diff != "" {
				t.Errorf("in: %v\n%s", in, diff)
			}
		}

		v0, err1 := n0.Value()
		v1, err2 := n1.Value()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}

		if diff := cmp.Diff(v0, v1); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}
	})
}

func FuzzString(f *testing.F) {
	f.Fuzz(func(t *testing.T, in string) {
		n0 := gnull.StringFrom(in)
		n1 := null.From[string](in)

		if diff := cmp.Diff(n0.IsZero(), n1.IsNull()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.ValueOrZero(), n1.ValueOrZero()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.Ptr(), n1.Ptr()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		b0, err1 := n0.MarshalJSON()
		b1, err2 := n1.MarshalJSON()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}
		if diff := cmp.Diff(string(b0), string(b1)); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		v0, err1 := n0.Value()
		v1, err2 := n1.Value()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}
		if diff := cmp.Diff(v0, v1); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}
	})
}

func FuzzTime(f *testing.F) {
	f.Fuzz(func(t *testing.T, sec int64, nsec int64, offset int64) {
		zone := time.FixedZone("", int(offset))
		in := time.Unix(sec, nsec).In(zone)
		n0 := gnull.TimeFrom(in)
		n1 := null.From[time.Time](in)

		if diff := cmp.Diff(n0.IsZero(), n1.IsNull()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.ValueOrZero(), n1.ValueOrZero()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		if diff := cmp.Diff(n0.Ptr(), n1.Ptr()); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		b0, err1 := n0.MarshalJSON()
		b1, err2 := n1.MarshalJSON()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}

		if diff := cmp.Diff(string(b0), string(b1)); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}

		v0, err1 := n0.Value()
		v1, err2 := n1.Value()
		if !errors.Is(err2, err1) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}
		if diff := cmp.Diff(v0, v1); diff != "" {
			t.Errorf("in: %v\n%s", in, diff)
		}
	})
}

func FuzzInt_UnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.Int
		err1 := json.Unmarshal(in, &n0)
		var n1 null.T[int64]
		err2 := json.Unmarshal(in, &n1)

		if err1 != nil && err2 == nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if err1 == nil && isIntJSONUnarshalIgnoreCase(t, in) {
			t.SkipNow()
		}

		if err1 == nil && err2 != nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.IsZero() != n1.IsNull() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzBool_UnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.Bool
		err1 := json.Unmarshal(in, &n0)
		var n1 null.T[bool]
		err2 := json.Unmarshal(in, &n1)

		if (err1 != nil) != (err2 != nil) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.IsZero() != n1.IsNull() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzFloat_UnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.Float
		err1 := json.Unmarshal(in, &n0)
		var n1 null.T[float64]
		err2 := json.Unmarshal(in, &n1)

		if err1 != nil && err2 == nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if err1 == nil && isFloatJSONUnmarshalIgnoreCase(t, in) {
			t.SkipNow()
		}

		if err1 == nil && err2 != nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.IsZero() != n1.IsNull() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzString_UnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.String
		err1 := json.Unmarshal(in, &n0)
		var n1 null.T[string]
		err2 := json.Unmarshal(in, &n1)

		if !errors.Is(err2, err1) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if err1 == nil && isFloatJSONUnmarshalIgnoreCase(t, in) {
			t.SkipNow()
		}

		if err1 == nil && err2 != nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.IsZero() != n1.IsNull() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzTime_UnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.Time
		err1 := json.Unmarshal(in, &n0)
		var n1 null.T[time.Time]
		err2 := json.Unmarshal(in, &n1)

		if !errors.Is(err2, err1) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if err1 == nil && isFloatJSONUnmarshalIgnoreCase(t, in) {
			t.SkipNow()
		}

		if err1 == nil && err2 != nil {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.IsZero() != n1.IsNull() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzInt_Scan_Int64(f *testing.F) {
	f.Fuzz(func(t *testing.T, in int64) {
		var n0 gnull.Int
		err1 := n0.Scan(in)
		var n1 null.T[int64]
		err2 := n1.Scan(in)

		// TODO: sentinel err
		if !errors.Is(err2, err1) {
			t.Errorf("in: %d, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %d, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzInt_Scan_Float64(f *testing.F) {
	f.Fuzz(func(t *testing.T, in float64) {
		var n0 gnull.Int
		err1 := n0.Scan(in)
		var n1 null.T[int64]
		err2 := n1.Scan(in)

		if (err1 != nil) != (err2 != nil) {
			t.Errorf("in: %v, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %v, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzInt_Scan_Bytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		var n0 gnull.Int
		err1 := n0.Scan(in)
		var n1 null.T[int64]
		err2 := n1.Scan(in)

		if !errors.Is(err2, err1) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzInt_Scan_String(f *testing.F) {
	f.Fuzz(func(t *testing.T, in string) {
		var n0 gnull.Int
		err1 := n0.Scan(in)
		var n1 null.T[int64]
		err2 := n1.Scan(in)

		if !errors.Is(err2, err1) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func FuzzInt_Scan_Time(f *testing.F) {
	f.Fuzz(func(t *testing.T, sec int64, nsec int64, offset int64) {
		zone := time.FixedZone("", int(offset))
		in := time.Unix(sec, nsec).In(zone)
		var n0 gnull.Int
		err1 := n0.Scan(in)
		var n1 null.T[int64]
		err2 := n1.Scan(in)

		if !errors.Is(err2, err1) {
			t.Errorf("in: %q, err1: %v, err2: %v", in, err1, err2)
		}

		if n0.ValueOrZero() != n1.ValueOrZero() {
			t.Errorf("in: %q, n0: %v, n1: %v", in, n0, n1)
		}
	})
}

func isFloatMarshalJSONIgnoreCase(t *testing.T, in float64) bool {
	const bits = strconv.IntSize
	abs := math.Abs(in)
	return bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21)
}

func isIntJSONUnarshalIgnoreCase(t *testing.T, in []byte) bool {
	b := in
	b = bytes.Trim(b, " \n\t\r")
	if len(b) < 3 {
		return false
	}
	if b[0] != '"' || b[len(b)-1] != '"' {
		return false
	}
	b = b[1 : len(b)-1]
	b, _ = bytes.CutPrefix(b, []byte("-"))
	b, _ = bytes.CutPrefix(b, []byte("+"))
	for _, c := range b {
		if c < '0' || '9' < c {
			return false
		}
	}
	return true
}

func isFloatJSONUnmarshalIgnoreCase(t *testing.T, in []byte) bool {
	var v float64
	err := json.Unmarshal(in, &v)
	if err == nil {
		return false
	}
	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		return false
	}
	if typeError.Value != "string" {
		return false
	}
	var str string
	if err := json.Unmarshal(in, &str); err != nil {
		return false
	}
	_, err = strconv.ParseFloat(str, 64)
	return err == nil
}
