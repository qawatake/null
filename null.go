package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	// can replace with "database/sql" after go1.22
	// https://github.com/golang/go/issues/60370
	sql1_22 "github.com/qawatake/null/internal/sql"
)

// MEMO: This package does not provide NewXXX functions.
// This is to ensure that when IsNull() returns true, ValueOrZero() guarantees to return the zero value.

// MEMO: This package does not provide SetValid function to make T immutable.

// MEMO: T implements neither MarshalText nor UnmarshalText.
// This is because there is no standard way to marshal/unmarshal for any type.

// T represents a value that may be null.
// The zero value for T is ready for use.
// It's not recommended to specify composite types with elements of reference types or interface types for V.
// This is because references may be unintentionally shared.
//
// If V is [strictly comparable], T[V] is also strictly comparable.
// However, if V implements method `Equal(u V) bool`, you should compare values of type T[V] using [T.Equal]. For details, refer to [T.Equal].
//
// [strictly comparable]: https://go.dev/ref/spec#Comparison_operators
type T[V comparable] struct {
	v sql1_22.Null[V]
}

// From creates a new T that is valid.
func From[V comparable](v V) T[V] {
	return T[V]{
		v: sql1_22.Null[V]{
			V:     v,
			Valid: true,
		},
	}
}

// FromPtr creates a new T that is null if p is nil.
func FromPtr[V comparable](p *V) T[V] {
	if p == nil {
		return T[V]{}
	}
	return From[V](*p)
}

var _ sql.Scanner = &T[int]{}

// Scan implements the sql.Scanner interface.
func (t *T[V]) Scan(src interface{}) error {
	if err := t.v.Scan(src); err != nil {
		*t = T[V]{}
		return err
	}
	if t.IsNull() {
		*t = T[V]{}
	}
	return nil
}

var _ driver.Valuer = T[int]{}

// Value implements the driver.Valuer interface.
func (t T[V]) Value() (driver.Value, error) {
	return t.v.Value()
}

var _ json.Unmarshaler = &T[int]{}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *T[V]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*t = T[V]{}
		return nil
	}
	var v V
	if err := json.Unmarshal(data, &v); err != nil {
		*t = T[V]{}
		return err
	}
	*t = From[V](v)
	return nil
}

var _ json.Marshaler = T[int]{}

// MarshalJSON implements the json.Marshaler interface.
func (t T[V]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("null"), nil
	}
	return json.Marshal(t.v.V)
}

var _ equaler[T[int]] = T[int]{}

type equaler[V any] interface {
	Equal(V) bool
}

// Equal reports whether t and u are equal.
// Two values t and u are equal if and only if either of the following conditions is met:
//   - t and u are both null.
//   - t and u are both not null and the internal values are equal in the sense of ==.
//   - t and u are both not null and the internal values are equal by `Equal(V) bool`.
//
// Even if t and u are different in terms of ==, they may be equal.
// So code should use Equal instead of == for comparison.
func (t T[V]) Equal(u T[V]) bool {
	if t.IsNull() && u.IsNull() {
		return true
	}
	if t.IsNull() && !u.IsNull() {
		return false
	}
	if !t.IsNull() && u.IsNull() {
		return false
	}
	if t == u {
		return true
	}
	if e, ok := any(t.v.V).(equaler[V]); ok {
		return e.Equal(u.v.V)
	}
	return false
}

// ValueOrZero returns the inner value V.
// If t is null (that is, t.IsNull() returns true), it returns the zero value of V.
func (t T[V]) ValueOrZero() V {
	if t.IsNull() {
		var v V
		return v
	}
	return t.v.V
}

// Ptr returns a pointer to the internal value, but it provides a different reference with each call.
// If t is null (that is, t.IsNull() returns true), it returns nil.
func (t T[V]) Ptr() *V {
	if t.IsNull() {
		return nil
	}
	v := t.v.V
	return &v
}

// IsNull reports whether t is null.
func (t T[V]) IsNull() bool {
	return !t.v.Valid
}

// nullBytes is a JSON null literal
var nullBytes = []byte("null")
