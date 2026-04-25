package dh

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const PostgresForeignKeyViolationCode = "23503"

type PostgresErrorDetails struct {
	Code           string
	ConstraintName string
}

func GetPostgresErrorDetails(err error) (PostgresErrorDetails, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return PostgresErrorDetails{}, false
	}

	return PostgresErrorDetails{
		Code:           pgErr.Code,
		ConstraintName: pgErr.ConstraintName,
	}, true
}

func IsPostgresForeignKeyViolation(err error) bool {
	details, ok := GetPostgresErrorDetails(err)
	if !ok {
		return false
	}

	return details.Code == PostgresForeignKeyViolationCode
}

func Reconcile[T comparable](incoming []T, current []T) (toDelete []T, toInsert []T) {
	incomingMap := make(map[T]struct{}, len(incoming))
	for _, id := range incoming {
		incomingMap[id] = struct{}{}
	}

	currentMap := make(map[T]struct{}, len(current))
	for _, id := range current {
		currentMap[id] = struct{}{}
	}

	for _, id := range current {
		if _, exists := incomingMap[id]; !exists {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range incoming {
		if _, exists := currentMap[id]; !exists {
			toInsert = append(toInsert, id)
		}
	}

	return
}
