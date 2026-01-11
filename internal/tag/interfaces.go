package tag

import (
	"context"
	"net/http"
)

// Repository defines the interface for tag data access
type Repository interface {
	GetByID(ctx context.Context, id string) (*Tag, error)
	GetByAlias(ctx context.Context, alias string) (*Tag, error)
	Create(ctx context.Context, tag *Tag) (string, error)
	IncrementUsageCount(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]Tag, error)
	Search(ctx context.Context, query string, limit int) ([]Tag, error)
}

// Service defines the interface for tag business logic
type Service interface {
	GetTag(ctx context.Context, id string) (*Tag, error)
	ListTags(ctx context.Context, limit, offset int) ([]Tag, error)
	SearchTags(ctx context.Context, query string, limit int) ([]Tag, error)
	NormalizeTags(ctx context.Context, rawTags []string) ([]string, error)
}

// Handler defines the interface for tag HTTP handlers
type Handler interface {
	GetTag(w http.ResponseWriter, r *http.Request)
	ListTags(w http.ResponseWriter, r *http.Request)
}
