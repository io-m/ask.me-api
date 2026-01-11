package feed

import (
	"context"
	"fmt"

	"github.com/askme/api/pkg/arango"
)

type repository struct {
	db *arango.Client
}

// NewRepository creates a new feed repository
func NewRepository(db *arango.Client) Repository {
	return &repository{db: db}
}

func (r *repository) GetRecommendedPosts(ctx context.Context, query FeedQuery) ([]FeedItem, string, error) {
	items, err := arango.Query[FeedItem](ctx, r.db, GetRecommendedPosts, map[string]any{
		"userId":   fmt.Sprintf("users/%s", query.UserID),
		"limit":    query.Limit,
		"category": query.Category,
		"depth":    query.Depth,
	})
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(items) == query.Limit && len(items) > 0 {
		nextCursor = fmt.Sprintf("%d", items[len(items)-1].CreatedAt)
	}

	return items, nextCursor, nil
}

func (r *repository) GetUserInteractionTags(ctx context.Context, userID string) ([]string, error) {
	return arango.Query[string](ctx, r.db, GetUserInteractionTags, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
	})
}

func (r *repository) GetUserCategories(ctx context.Context, userID string) ([]string, error) {
	return arango.Query[string](ctx, r.db, GetUserCategories, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
	})
}

func (r *repository) GetUserIntents(ctx context.Context, userID string) ([]string, error) {
	return arango.Query[string](ctx, r.db, GetUserIntents, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
	})
}
