package feed

import (
	"context"
	"net/http"
)

// Repository defines the interface for feed data access
type Repository interface {
	GetRecommendedPosts(ctx context.Context, query FeedQuery) ([]FeedItem, string, error)
	GetUserInteractionTags(ctx context.Context, userID string) ([]string, error)
	GetUserCategories(ctx context.Context, userID string) ([]string, error)
	GetUserIntents(ctx context.Context, userID string) ([]string, error)
}

// Service defines the interface for feed business logic
type Service interface {
	GetFeed(ctx context.Context, query FeedQuery) (*FeedResponse, error)
}

// Handler defines the interface for feed HTTP handlers
type Handler interface {
	GetFeed(w http.ResponseWriter, r *http.Request)
}
