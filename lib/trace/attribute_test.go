package trace

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestAttribute_Slog(t *testing.T) {
	tests := []struct {
		name     string
		attr     Attribute
		wantName string
	}{
		{
			name:     "string value",
			attr:     Attribute{Name: "test", Value: "value"},
			wantName: "test",
		},
		{
			name:     "int value",
			attr:     Attribute{Name: "test_int", Value: 123},
			wantName: "test_int",
		},
		{
			name:     "bool value",
			attr:     Attribute{Name: "test_bool", Value: true},
			wantName: "test_bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slogAttr := tt.attr.Slog()
			assert.Equal(t, tt.wantName, slogAttr.Key)
			// 元の値とAny()の型が若干異なることがあるため、
			// 厳密な型比較を避け、文字列表現で比較
			assert.Equal(t, fmt.Sprintf("%v", tt.attr.Value), fmt.Sprintf("%v", slogAttr.Value.Any()))
		})
	}
}

func TestAttribute_Trace(t *testing.T) {
	tests := []struct {
		name     string
		attr     Attribute
		wantName string
	}{
		{
			name:     "string value",
			attr:     Attribute{Name: "test", Value: "value"},
			wantName: "test",
		},
		{
			name:     "int value",
			attr:     Attribute{Name: "test_int", Value: 123},
			wantName: "test_int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traceAttr := tt.attr.Trace()
			assert.Equal(t, attribute.Key(tt.wantName), traceAttr.Key)
		})
	}
}

type mockStringer struct {
	val string
}

func (m mockStringer) String() string {
	return m.val
}

func TestAnyValue(t *testing.T) {
	now := time.Now()
	duration := 10 * time.Second

	tests := []struct {
		name  string
		value interface{}
		check func(*testing.T, attribute.Value)
	}{
		{
			name:  "string",
			value: "test string",
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, "test string", v.AsString())
			},
		},
		{
			name:  "[]string",
			value: []string{"a", "b", "c"},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, []string{"a", "b", "c"}, v.AsStringSlice())
			},
		},
		{
			name:  "int",
			value: 42,
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, int64(42), v.AsInt64())
			},
		},
		{
			name:  "[]int",
			value: []int{1, 2, 3},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, []int64{1, 2, 3}, v.AsInt64Slice())
			},
		},
		{
			name:  "bool",
			value: true,
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, true, v.AsBool())
			},
		},
		{
			name:  "[]bool",
			value: []bool{true, false, true},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, []bool{true, false, true}, v.AsBoolSlice())
			},
		},
		{
			name:  "time.Duration",
			value: duration,
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, int64(duration), v.AsInt64())
			},
		},
		{
			name:  "time.Time",
			value: now,
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, now.Format(time.RFC3339), v.AsString())
			},
		},
		{
			name:  "float64",
			value: 3.14,
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, 3.14, v.AsFloat64())
			},
		},
		{
			name:  "[]float64",
			value: []float64{1.1, 2.2, 3.3},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, []float64{1.1, 2.2, 3.3}, v.AsFloat64Slice())
			},
		},
		{
			name:  "float32",
			value: float32(2.5),
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, 2.5, v.AsFloat64())
			},
		},
		{
			name:  "fmt.Stringer",
			value: mockStringer{val: "mock value"},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, "mock value", v.AsString())
			},
		},
		{
			name:  "[]fmt.Stringer",
			value: []fmt.Stringer{mockStringer{val: "m1"}, mockStringer{val: "m2"}},
			check: func(t *testing.T, v attribute.Value) {
				assert.Equal(t, []string{"m1", "m2"}, v.AsStringSlice())
			},
		},
		{
			name:  "default",
			value: struct{ Name string }{Name: "test"},
			check: func(t *testing.T, v attribute.Value) {
				assert.Contains(t, v.AsString(), "{test}")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := anyValue(tt.value)
			tt.check(t, val)
		})
	}
}

func TestAttributes_Slog(t *testing.T) {
	attrs := Attributes{
		{Name: "str", Value: "value"},
		{Name: "num", Value: 123},
		{Name: "bool", Value: true},
	}

	results := attrs.Slog()
	assert.Len(t, results, 3)
}

func TestAttributes_Trace(t *testing.T) {
	attrs := Attributes{
		{Name: "str", Value: "value"},
		{Name: "num", Value: 123},
		{Name: "bool", Value: true},
	}

	results := attrs.Trace()
	assert.Len(t, results, 3)
	assert.Equal(t, attribute.Key("str"), results[0].Key)
	assert.Equal(t, attribute.Key("num"), results[1].Key)
	assert.Equal(t, attribute.Key("bool"), results[2].Key)
}
