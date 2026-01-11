package post

import (
	"context"
	"net/http"
)

// Repository defines the interface for post data access
type Repository interface {
	GetByID(ctx context.Context, id string) (*Post, error)
	Create(ctx context.Context, post *Post) (string, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error

	// Edge operations
	CreateCreatedEdge(ctx context.Context, userID, postID string, createdAt int64) error
	CreateRespondedEdge(ctx context.Context, userID, postID, chatID string, createdAt int64) error
	CreateVotedEdge(ctx context.Context, userID, postID, option string, createdAt int64) error
	CreatePostHasTagEdge(ctx context.Context, postID, tagKey string, confidence float64) error

	// Query operations
	GetPostTags(ctx context.Context, postID string) ([]string, error)
	GetVotes(ctx context.Context, postID string) (map[string]int, error)
	HasUserVoted(ctx context.Context, userID, postID string) (bool, error)
	HasUserResponded(ctx context.Context, userID, postID string) (bool, error)
	GetAuthor(ctx context.Context, postID string) (*PostAuthor, error)
}

// Service defines the interface for post business logic
type Service interface {
	GetPost(ctx context.Context, id string) (*GetPostResponse, error)
	CreatePost(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error)
	CreatePoll(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error)
	RespondToPost(ctx context.Context, postID string, req *RespondToPostRequest) (*RespondToPostResponse, error)
	Vote(ctx context.Context, postID string, req *VoteRequest) (*VoteResponse, error)
}

// Handler defines the interface for post HTTP handlers
type Handler interface {
	GetPost(w http.ResponseWriter, r *http.Request)
	CreatePost(w http.ResponseWriter, r *http.Request)
	CreatePoll(w http.ResponseWriter, r *http.Request)
	RespondToPost(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
}
