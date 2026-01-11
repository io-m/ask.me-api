package tag

import (
	"context"

	"github.com/askme/api/pkg/arango"
)

type repository struct {
	db *arango.Client
}

// NewRepository creates a new tag repository
func NewRepository(db *arango.Client) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*Tag, error) {
	return arango.QueryOne[Tag](ctx, r.db, GetTagByID, map[string]any{"key": id})
}

func (r *repository) GetByAlias(ctx context.Context, alias string) (*Tag, error) {
	return arango.QueryOne[Tag](ctx, r.db, GetTagByAlias, map[string]any{"alias": alias})
}

func (r *repository) Create(ctx context.Context, tag *Tag) (string, error) {
	return arango.InsertDocument(ctx, r.db, arango.CollectionTags, tag)
}

func (r *repository) IncrementUsageCount(ctx context.Context, id string) error {
	_, err := arango.Query[any](ctx, r.db, IncrementTagUsageCount, map[string]any{"key": id})
	return err
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]Tag, error) {
	return arango.Query[Tag](ctx, r.db, ListTagsByUsage, map[string]any{
		"limit":  limit,
		"offset": offset,
	})
}

func (r *repository) Search(ctx context.Context, searchQuery string, limit int) ([]Tag, error) {
	return arango.Query[Tag](ctx, r.db, SearchTags, map[string]any{
		"query": searchQuery,
		"limit": limit,
	})
}
