//go:build integration
// +build integration

package testsupport

import (
	"context"
	"log/slog"
	"testing"

	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type Fixtures struct {
	Context context.Context
	DB      *sqlx.DB
	Logger  *slog.Logger
}

func NewFixtures(ctx context.Context, db *sqlx.DB, logger *slog.Logger) Fixtures {
	return Fixtures{
		Context: ctx,
		DB:      db,
		Logger:  logger,
	}
}

func (f Fixtures) MustCreateItem(t testing.TB, create *domains.ItemCreate) uuid.UUID {
	t.Helper()

	repo := repositories.NewItemsRepository(f.DB, f.Logger)
	id, err := repo.Create(f.Context, create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	return id
}

func (f Fixtures) MustCreateTag(t testing.TB, create *domains.TagCreate) uuid.UUID {
	t.Helper()

	repo := repositories.NewTagsRepository(f.DB, f.Logger)
	id, err := repo.Create(f.Context, create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	return id
}

func (f Fixtures) MustLinkItemTags(t testing.TB, itemID uuid.UUID, tagIDs ...uuid.UUID) {
	t.Helper()

	repo := repositories.NewTagToItemsRepository(f.DB, f.Logger)
	tagRefs := make([]*uuid.UUID, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tagIDCopy := tagID
		tagRefs = append(tagRefs, &tagIDCopy)
	}

	ok, err := repo.BulkInsert(f.Context, &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: tagRefs,
	})
	require.NoError(t, err)
	require.True(t, ok)
}
