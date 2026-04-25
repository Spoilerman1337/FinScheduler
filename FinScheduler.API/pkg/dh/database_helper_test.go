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
	t.Run("extracts details from wrapped pg error", func(t *testing.T) {
		// Arrange
		err := fmt.Errorf("wrapped: %w", &pgconn.PgError{
			Code:           PostgresForeignKeyViolationCode,
			ConstraintName: "fk_tag_to_item_tag_id",
		})

		// Act
		details, ok := GetPostgresErrorDetails(err)

		// Assert
		require.True(t, ok)
		assert.Equal(t, PostgresForeignKeyViolationCode, details.Code)
		assert.Equal(t, "fk_tag_to_item_tag_id", details.ConstraintName)
	})

	t.Run("returns false for non pg error", func(t *testing.T) {
		// Arrange

		// Act
		details, ok := GetPostgresErrorDetails(errors.New("boom"))

		// Assert
		require.False(t, ok)
		assert.Equal(t, PostgresErrorDetails{}, details)
	})
}

func TestIsPostgresForeignKeyViolation(t *testing.T) {
	t.Run("returns true for foreign key violation", func(t *testing.T) {
		// Arrange
		err := &pgconn.PgError{Code: PostgresForeignKeyViolationCode}

		// Act
		isViolation := IsPostgresForeignKeyViolation(err)

		// Assert
		assert.True(t, isViolation)
	})

	t.Run("returns false for other postgres errors", func(t *testing.T) {
		// Arrange
		err := &pgconn.PgError{Code: "23505"}

		// Act
		isViolation := IsPostgresForeignKeyViolation(err)

		// Assert
		assert.False(t, isViolation)
	})

	t.Run("returns false for non postgres errors", func(t *testing.T) {
		// Arrange
		err := errors.New("boom")

		// Act
		isViolation := IsPostgresForeignKeyViolation(err)

		// Assert
		assert.False(t, isViolation)
	})
}

func TestReconcile(t *testing.T) {
	t.Run("returns inserts and deletes for changed sets", func(t *testing.T) {
		// Arrange

		// Act
		toDelete, toInsert := Reconcile([]int{2, 3}, []int{1, 2})

		// Assert
		assert.Equal(t, []int{1}, toDelete)
		assert.Equal(t, []int{3}, toInsert)
	})

	t.Run("returns empty slices for identical sets", func(t *testing.T) {
		// Arrange

		// Act
		toDelete, toInsert := Reconcile([]int{1, 2}, []int{1, 2})

		// Assert
		assert.Empty(t, toDelete)
		assert.Empty(t, toInsert)
	})

	t.Run("returns all incoming values for insertion when current is empty", func(t *testing.T) {
		// Arrange

		// Act
		toDelete, toInsert := Reconcile([]int{1, 2}, []int{})

		// Assert
		assert.Empty(t, toDelete)
		assert.Equal(t, []int{1, 2}, toInsert)
	})
}
