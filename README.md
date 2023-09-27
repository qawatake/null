# null

[![test](https://github.com/qawatake/null/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/null/actions/workflows/test.yaml)

`null` is a package for providing types to handle nullable values.

For a `comparable` type `V`, `null.T[V]` implements methods like `sql.Scanner` and `json.Unmarshaler`. For instance, it supports types such as:

- Derived integer types like `time.Duration`
- Structures with `comparable` fields
- Arrays consisting of `comparable` elements

`null.T` is oriented towards *immutability*. Therefore, unlike `sql.NullInt64`, `null.T` does not have APIs for modification.

- `null.T` does not expose its fields.
- `null.T` does not have methods for modification (excluding `Scan` and `Unmarshal`).

```go
func main() {
	var d null.T[time.Duration]
	json.Unmarshal([]byte(`0`), &d)
	fmt.Printf("valid: %v, value: %v\n", !d.IsNull(), d.ValueOrZero())
	// Output:
	// valid: true, value: 0s

	json.Unmarshal([]byte(`null`), &d)
	fmt.Printf("valid: %v, value: %v\n", !d.IsNull(), d.ValueOrZero())
	// Output:
	// valid: false, value: 0s
}
```

## Differences from [gopkg.in/guregu/null]

Differences from the well-known package [gopkg.in/guregu/null], which also defines nullable types include:

- Generics: `null` supports not only `bool`, `int64`, `float64`, `string`, and `time.Time`, but any `comparable` type.
- Immutable: `null.T` does not have APIs for modification.
  - `null.T` does not expose its fields.
  - `null.T` does not have a `SetValid` method.
- Fewer APIs: This package does not provide several APIs that are present in guregu/null.
  - `New*`
  - `MarshalText`, `UnmarshalText`
  - `SetValid`

### Minor differences

- `guregu/null.Int` can json unmarshal string types like `"123"`, but `null.T[int64]` will result in an error. This aligns with the behavior when `json.Unmarshal` is applied to `int64`.
- `guregu/null.Float` can json unmarshal string types like `"1.23"`, but `null.T[float64]` will result in an error. This aligns with the behavior when `json.Unmarshal` is applied to `float64`.
- For floats with extremely small or large absolute values, the output of `MarshalJSON` for `guregu/null.Float` and `null.T[float64]` is different. `null.T[float64]` matches the behavior when `float64` is passed to `json.MarshalJSON`.
- The method `IsZero` has been renamed to `IsNull`.

<!-- links -->
[gopkg.in/guregu/null]: https://github.com/guregu/null
