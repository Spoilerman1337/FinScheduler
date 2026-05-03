package rh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDereferenceSlice(t *testing.T) {
	firstValue := 10
	secondValue := 20

	tests := []struct {
		name     string
		input    []*int
		expected []int
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: []int{},
		},
		{
			name: "all pointers are non nil",
			input: []*int{
				&firstValue,
				&secondValue,
			},
			expected: []int{10, 20},
		},
		{
			name: "nil pointers become zero values",
			input: []*int{
				&firstValue,
				nil,
				&secondValue,
			},
			expected: []int{10, 0, 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			input := tt.input

			// Act
			result := DereferenceSlice(input)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReferenceSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: []int{},
		},
		{
			name:     "values are referenced in order",
			input:    []int{10, 20, 30},
			expected: []int{10, 20, 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			input := tt.input

			// Act
			result := ReferenceSlice(input)

			// Assert
			require.Len(t, result, len(tt.expected))

			for index := range tt.expected {
				require.NotNil(t, result[index])
				assert.Equal(t, tt.expected[index], *result[index])
			}
		})
	}
}
