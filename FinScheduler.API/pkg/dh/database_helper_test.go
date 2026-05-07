package dh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPostgresErrorDetails(t *testing.T) {
	tests := []struct {
		name             string
		inputError       error
		expectedDetails  PostgresErrorDetails
		expectedDetected bool
	}{
		{
			name: "postgres error",
			inputError: &pgconn.PgError{
				Code:           PostgresForeignKeyViolationCode,
				ConstraintName: "tag_to_item_tag_id_fkey",
			},
			expectedDetails: PostgresErrorDetails{
				Code:           PostgresForeignKeyViolationCode,
				ConstraintName: "tag_to_item_tag_id_fkey",
			},
			expectedDetected: true,
		},
		{
			name: "wrapped postgres error",
			inputError: fmt.Errorf("wrapped: %w", &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "items_name_key",
			}),
			expectedDetails: PostgresErrorDetails{
				Code:           "23505",
				ConstraintName: "items_name_key",
			},
			expectedDetected: true,
		},
		{
			name:             "non postgres error",
			inputError:       errors.New("plain error"),
			expectedDetails:  PostgresErrorDetails{},
			expectedDetected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			inputError := tt.inputError

			// Act
			details, ok := GetPostgresErrorDetails(inputError)

			// Assert
			assert.Equal(t, tt.expectedDetected, ok)
			assert.Equal(t, tt.expectedDetails, details)
		})
	}
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name             string
		incoming         []int
		current          []int
		expectedToDelete []int
		expectedToInsert []int
	}{
		{
			name:             "same sets",
			incoming:         []int{1, 2},
			current:          []int{1, 2},
			expectedToDelete: nil,
			expectedToInsert: nil,
		},
		{
			name:             "incoming has extra values",
			incoming:         []int{1, 2, 3},
			current:          []int{1, 2},
			expectedToDelete: nil,
			expectedToInsert: []int{3},
		},
		{
			name:             "current has extra values",
			incoming:         []int{1},
			current:          []int{1, 2},
			expectedToDelete: []int{2},
			expectedToInsert: nil,
		},
		{
			name:             "partial overlap",
			incoming:         []int{1, 3},
			current:          []int{1, 2},
			expectedToDelete: []int{2},
			expectedToInsert: []int{3},
		},
		{
			name:             "duplicates are preserved in result order",
			incoming:         []int{1, 3, 3},
			current:          []int{1, 2, 2},
			expectedToDelete: []int{2, 2},
			expectedToInsert: []int{3, 3},
		},
		{
			name:             "empty slices",
			incoming:         nil,
			current:          nil,
			expectedToDelete: nil,
			expectedToInsert: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			incoming := tt.incoming
			current := tt.current

			// Act
			toDelete, toInsert := Reconcile(incoming, current)

			// Assert
			require.Equal(t, tt.expectedToDelete, toDelete)
			require.Equal(t, tt.expectedToInsert, toInsert)
		})
	}
}
