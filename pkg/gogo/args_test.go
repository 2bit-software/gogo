package gogo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fieldOf creates a settable reflect.Value of the given type for testing.
func fieldOf[T any]() reflect.Value {
	var v T
	return reflect.ValueOf(&v).Elem()
}

func TestSetFieldFromString(t *testing.T) {
	tests := []struct {
		name      string
		field     reflect.Value
		value     string
		fieldName string
		wantErr   string
		check     func(t *testing.T, field reflect.Value)
	}{
		// String
		{
			name:      "string set hello",
			field:     fieldOf[string](),
			value:     "hello",
			fieldName: "Name",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, "hello", f.String())
			},
		},
		{
			name:      "string set empty",
			field:     fieldOf[string](),
			value:     "",
			fieldName: "Name",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, "", f.String())
			},
		},

		// Bool true variants
		{name: "bool true", field: fieldOf[bool](), value: "true", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},
		{name: "bool 1", field: fieldOf[bool](), value: "1", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},
		{name: "bool t", field: fieldOf[bool](), value: "t", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},
		{name: "bool TRUE", field: fieldOf[bool](), value: "TRUE", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},
		{name: "bool True", field: fieldOf[bool](), value: "True", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},
		{name: "bool T", field: fieldOf[bool](), value: "T", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.True(t, f.Bool()) }},

		// Bool false variants
		{name: "bool false", field: fieldOf[bool](), value: "false", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},
		{name: "bool 0", field: fieldOf[bool](), value: "0", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},
		{name: "bool f", field: fieldOf[bool](), value: "f", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},
		{name: "bool FALSE", field: fieldOf[bool](), value: "FALSE", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},
		{name: "bool False", field: fieldOf[bool](), value: "False", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},
		{name: "bool F", field: fieldOf[bool](), value: "F", fieldName: "Flag",
			check: func(t *testing.T, f reflect.Value) { assert.False(t, f.Bool()) }},

		// Bool invalid
		{name: "bool yes invalid", field: fieldOf[bool](), value: "yes", fieldName: "Flag",
			wantErr: "invalid boolean value for"},
		{name: "bool no invalid", field: fieldOf[bool](), value: "no", fieldName: "Flag",
			wantErr: "invalid boolean value for"},
		{name: "bool maybe invalid", field: fieldOf[bool](), value: "maybe", fieldName: "Flag",
			wantErr: "invalid boolean value for"},
		{name: "bool empty invalid", field: fieldOf[bool](), value: "", fieldName: "Flag",
			wantErr: "invalid boolean value for"},

		// Int
		{
			name:      "int 42",
			field:     fieldOf[int](),
			value:     "42",
			fieldName: "Count",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, int64(42), f.Int())
			},
		},
		{name: "int abc invalid", field: fieldOf[int](), value: "abc", fieldName: "Count",
			wantErr: "invalid integer value for"},

		// Int8 boundaries
		{
			name:      "int8 max 127",
			field:     fieldOf[int8](),
			value:     "127",
			fieldName: "Small",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, int64(127), f.Int())
			},
		},
		{name: "int8 overflow 128", field: fieldOf[int8](), value: "128", fieldName: "Small",
			wantErr: "overflows"},
		{
			name:      "int8 min -128",
			field:     fieldOf[int8](),
			value:     "-128",
			fieldName: "Small",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, int64(-128), f.Int())
			},
		},
		{name: "int8 underflow -129", field: fieldOf[int8](), value: "-129", fieldName: "Small",
			wantErr: "overflows"},

		// Uint8 boundaries
		{
			name:      "uint8 max 255",
			field:     fieldOf[uint8](),
			value:     "255",
			fieldName: "Byte",
			check: func(t *testing.T, f reflect.Value) {
				assert.Equal(t, uint64(255), f.Uint())
			},
		},
		{name: "uint8 overflow 256", field: fieldOf[uint8](), value: "256", fieldName: "Byte",
			wantErr: "overflows"},
		{name: "uint8 negative -1", field: fieldOf[uint8](), value: "-1", fieldName: "Byte",
			wantErr: "invalid unsigned integer value for"},

		// Float64
		{
			name:      "float64 3.14",
			field:     fieldOf[float64](),
			value:     "3.14",
			fieldName: "Rate",
			check: func(t *testing.T, f reflect.Value) {
				assert.InDelta(t, 3.14, f.Float(), 0.001)
			},
		},
		{name: "float64 abc invalid", field: fieldOf[float64](), value: "abc", fieldName: "Rate",
			wantErr: "invalid float value for"},

		// Float32
		{
			name:      "float32 3.14",
			field:     fieldOf[float32](),
			value:     "3.14",
			fieldName: "SmallRate",
			check: func(t *testing.T, f reflect.Value) {
				assert.InDelta(t, 3.14, f.Float(), 0.01)
			},
		},
		{
			name:      "float32 overflow",
			field:     fieldOf[float32](),
			value:     "3.5e38",
			fieldName: "SmallRate",
			wantErr:   "overflows",
		},

		// Unsupported type
		{
			name:      "unsupported slice type",
			field:     fieldOf[[]string](),
			value:     "anything",
			fieldName: "Tags",
			wantErr:   "unsupported field type for",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setFieldFromString(tt.field, tt.value, tt.fieldName)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, tt.field)
				}
			}
		})
	}
}

func TestHydrateFromPositional(t *testing.T) {
	t.Run("basic assignment", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"1"`
			C string `order:"2"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x", "y", "z"})
		assert.NoError(t, err)
		assert.Equal(t, "x", o.A)
		assert.Equal(t, "y", o.B)
		assert.Equal(t, "z", o.C)
	})

	t.Run("pre-set field is skipped", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"1"`
		}
		o := &opts{A: "existing"}
		err := HydrateFromPositional(o, []string{"x", "y"})
		assert.NoError(t, err)
		assert.Equal(t, "existing", o.A)
		assert.Equal(t, "y", o.B)
	})

	t.Run("empty quotes consumed but field not set", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"1"`
		}
		// 2-byte string: two literal double-quote characters
		o := &opts{}
		err := HydrateFromPositional(o, []string{`""`, "hello"})
		assert.NoError(t, err)
		assert.Equal(t, "", o.A)
		assert.Equal(t, "hello", o.B)
	})

	t.Run("escaped quotes consumed but field not set", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"1"`
		}
		// 4-byte string: backslash-quote-backslash-quote
		o := &opts{}
		err := HydrateFromPositional(o, []string{`\"\"`, "hello"})
		assert.NoError(t, err)
		assert.Equal(t, "", o.A)
		assert.Equal(t, "hello", o.B)
	})

	t.Run("fewer args than fields", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"1"`
			C string `order:"2"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x"})
		assert.NoError(t, err)
		assert.Equal(t, "x", o.A)
		assert.Equal(t, "", o.B)
		assert.Equal(t, "", o.C)
	})

	t.Run("non-pointer input", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
		}
		err := HydrateFromPositional(opts{}, []string{"x"})
		assert.ErrorContains(t, err, "expected pointer to struct")
	})

	t.Run("pointer to non-struct", func(t *testing.T) {
		s := "hello"
		err := HydrateFromPositional(&s, []string{"x"})
		assert.ErrorContains(t, err, "expected pointer to struct")
	})

	t.Run("nil input", func(t *testing.T) {
		err := HydrateFromPositional(nil, []string{"x"})
		assert.Error(t, err)
	})

	t.Run("invalid order tag", func(t *testing.T) {
		type opts struct {
			A string `order:"abc"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x"})
		assert.ErrorContains(t, err, "invalid order tag for field")
	})

	t.Run("unexported fields skipped", func(t *testing.T) {
		type opts struct {
			hidden string `order:"0"`
			Public string `order:"1"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x", "y"})
		assert.NoError(t, err)
		assert.Equal(t, "", o.hidden)
		assert.Equal(t, "y", o.Public)
	})

	t.Run("more args than fields", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x", "y", "z"})
		assert.NoError(t, err)
		assert.Equal(t, "x", o.A)
	})

	t.Run("empty args", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{})
		assert.NoError(t, err)
		assert.Equal(t, "", o.A)
	})

	t.Run("no order tags", func(t *testing.T) {
		type opts struct {
			A string
			B string
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"x", "y"})
		assert.NoError(t, err)
		assert.Equal(t, "", o.A)
		assert.Equal(t, "", o.B)
	})

	t.Run("gapped order tags", func(t *testing.T) {
		type opts struct {
			A string `order:"0"`
			B string `order:"3"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"first", "skip1", "skip2", "fourth"})
		assert.NoError(t, err)
		assert.Equal(t, "first", o.A)
		assert.Equal(t, "fourth", o.B)
	})

	t.Run("mixed types with conversion", func(t *testing.T) {
		type opts struct {
			Name  string `order:"0"`
			Count int    `order:"1"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"hello", "42"})
		assert.NoError(t, err)
		assert.Equal(t, "hello", o.Name)
		assert.Equal(t, 42, o.Count)
	})

	t.Run("type conversion error propagated", func(t *testing.T) {
		type opts struct {
			Val float32 `order:"0"`
		}
		o := &opts{}
		err := HydrateFromPositional(o, []string{"3.5e38"})
		assert.ErrorContains(t, err, "overflows")
	})
}

func TestParseArgs(t *testing.T) {
	t.Run("flag sets field no positional", func(t *testing.T) {
		type opts struct {
			Name string `long:"name"`
		}
		o := &opts{}
		positional, err := ParseArgs(o, []string{"--name", "foo"})
		assert.NoError(t, err)
		assert.Equal(t, "foo", o.Name)
		assert.Empty(t, positional)
	})

	t.Run("flag plus positional args", func(t *testing.T) {
		type opts struct {
			Name string `long:"name"`
		}
		o := &opts{}
		positional, err := ParseArgs(o, []string{"--name", "foo", "bar"})
		assert.NoError(t, err)
		assert.Equal(t, "foo", o.Name)
		assert.Equal(t, []string{"bar"}, positional)
	})

	t.Run("unknown flag returns error", func(t *testing.T) {
		type opts struct {
			Name string `long:"name"`
		}
		o := &opts{}
		_, err := ParseArgs(o, []string{"--unknown"})
		assert.Error(t, err)
	})
}
