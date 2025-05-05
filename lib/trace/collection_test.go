package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapFunc  func(int) string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []int{},
			mapFunc:  func(i int) string { return string(rune(i + 'a')) },
			expected: []string{},
		},
		{
			name:     "transform ints to strings",
			input:    []int{1, 2, 3},
			mapFunc:  func(i int) string { return string(rune(i + 'a')) },
			expected: []string{"b", "c", "d"},
		},
		{
			name:     "transform with custom function",
			input:    []int{10, 20, 30},
			mapFunc:  func(i int) string { return string(rune(i/10 + '0')) },
			expected: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := arrayMap(tt.input, tt.mapFunc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt64Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []int64
	}{
		{
			name:     "[]int",
			input:    []int{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]int8",
			input:    []int8{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]int16",
			input:    []int16{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]int32",
			input:    []int32{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]int64",
			input:    []int64{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]uint",
			input:    []uint{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			name:     "[]uint8",
			input:    []uint8{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case []int:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []int8:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []int16:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []int32:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []int64:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []uint:
				assert.Equal(t, tt.expected, int64Slice(v))
			case []uint8:
				assert.Equal(t, tt.expected, int64Slice(v))
			}
		})
	}
}

func TestFloat32ToFloat64Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected []float64
	}{
		{
			name:     "empty slice",
			input:    []float32{},
			expected: []float64{},
		},
		{
			name:     "basic conversion",
			input:    []float32{1.1, 2.2, 3.3},
			expected: []float64{1.1, 2.2, 3.3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float32ToFloat64Slice(tt.input)

			// Float32をFloat64に変換するとわずかな精度の差が生じる可能性があるため
			// 厳密な等価をチェックする代わりに、各要素を個別に比較
			assert.Equal(t, len(tt.expected), len(result))
			for i := range result {
				assert.InDelta(t, tt.expected[i], result[i], 0.00001)
			}
		})
	}
}
