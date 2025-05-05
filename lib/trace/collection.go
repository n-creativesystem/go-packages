package trace

func arrayMap[T any, V any](values []T, fn func(value T) V) []V {
	results := make([]V, len(values))
	for idx, value := range values {
		results[idx] = fn(value)
	}
	return results
}

func int64Slice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr](values []T) []int64 {
	return arrayMap(values, func(v T) int64 { return int64(v) })
}

func float32ToFloat64Slice(values []float32) []float64 {
	results := make([]float64, len(values))
	for idx, value := range values {
		results[idx] = float64(value)
	}
	return results
}
