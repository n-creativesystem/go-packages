package trace

import (
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type Attribute struct {
	Name  string
	Value any
}

func (a *Attribute) Slog() slog.Attr {
	return slog.Any(a.Name, a.Value)
}

func (a *Attribute) Trace() attribute.KeyValue {
	return attribute.KeyValue{Key: attribute.Key(a.Name), Value: anyValue(a.Value)}
}

func anyValue(v any) attribute.Value {
	switch v := v.(type) {
	case string:
		return attribute.StringValue(v)
	case []string:
		return attribute.StringSliceValue(v)
	case int:
		return attribute.Int64Value(int64(v))
	case []int:
		return attribute.Int64SliceValue(int64Slice(v))
	case uint:
		return attribute.Int64Value(int64(v))
	case []uint:
		return attribute.Int64SliceValue(int64Slice(v))
	case int64:
		return attribute.Int64Value(v)
	case []int64:
		return attribute.Int64SliceValue(int64Slice(v))
	case uint64:
		return attribute.Int64Value(int64(v))
	case []uint64:
		return attribute.Int64SliceValue(int64Slice(v))
	case bool:
		return attribute.BoolValue(v)
	case []bool:
		return attribute.BoolSliceValue(v)
	case time.Duration:
		return attribute.Int64Value(int64(v))
	case time.Time:
		return attribute.StringValue(v.Format(time.RFC3339))
	case uint8:
		return attribute.Int64Value(int64(v))
	case []uint8:
		return attribute.Int64SliceValue(int64Slice(v))
	case uint16:
		return attribute.Int64Value(int64(v))
	case []uint16:
		return attribute.Int64SliceValue(int64Slice(v))
	case uint32:
		return attribute.Int64Value(int64(v))
	case []uint32:
		return attribute.Int64SliceValue(int64Slice(v))
	case uintptr:
		return attribute.Int64Value(int64(v))
	case []uintptr:
		return attribute.Int64SliceValue(int64Slice(v))
	case int8:
		return attribute.Int64Value(int64(v))
	case []int8:
		return attribute.Int64SliceValue(int64Slice(v))
	case int16:
		return attribute.Int64Value(int64(v))
	case []int16:
		return attribute.Int64SliceValue(int64Slice(v))
	case int32:
		return attribute.Int64Value(int64(v))
	case []int32:
		return attribute.Int64SliceValue(int64Slice(v))
	case float64:
		return attribute.Float64Value(v)
	case []float64:
		return attribute.Float64SliceValue(v)
	case float32:
		return attribute.Float64Value(float64(v))
	case []float32:
		return attribute.Float64SliceValue(float32ToFloat64Slice(v))
	case fmt.Stringer:
		return attribute.StringValue(v.String())
	case []fmt.Stringer:
		return attribute.StringSliceValue(arrayMap(v, func(value fmt.Stringer) string { return value.String() }))
	default:
		return attribute.StringValue(fmt.Sprintf("%v", v))
	}
}

type Attributes []Attribute

func (attrs Attributes) Slog() []any {
	results := make([]any, len(attrs))
	for idx, attr := range attrs {
		results[idx] = attr.Slog()
	}
	return results
}

func (attrs Attributes) Trace() []attribute.KeyValue {
	results := make([]attribute.KeyValue, len(attrs))
	for idx, attr := range attrs {
		results[idx] = attr.Trace()
	}
	return results
}
