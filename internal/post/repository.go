package post

import (
	"context"
	"fmt"

	"github.com/askme/api/pkg/arango"
)

type repository struct {
	db *arango.Client
}

// NewRepository creates a new post repository
func NewRepository(db *arango.Client) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*Post, error) {
	return arango.QueryOne[Post](ctx, r.db, GetPostByID, map[string]any{"key": id})
}

func (r *repository) Create(ctx context.Context, post *Post) (string, error) {
	return arango.InsertDocument(ctx, r.db, arango.CollectionPosts, post)
}

func (r *repository) Update(ctx context.Context, post *Post) error {
	return arango.UpdateDocument(ctx, r.db, arango.CollectionPosts, post.Key, post)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return arango.DeleteDocument(ctx, r.db, arango.CollectionPosts, id)
}

func (r *repository) CreateCreatedEdge(ctx context.Context, userID, postID string, createdAt int64) error {
	edge := CreatedEdge{
		From:      fmt.Sprintf("users/%s", userID),
		To:        fmt.Sprintf("posts/%s", postID),
		CreatedAt: createdAt,
	}
	_, err := arango.InsertDocument(ctx, r.db, arango.EdgeCreated, edge)
	return err
}

func (r *repository) CreateRespondedEdge(ctx context.Context, userID, postID, chatID string, createdAt int64) error {
	edge := RespondedEdge{
		From:      fmt.Sprintf("users/%s", userID),
		To:        fmt.Sprintf("posts/%s", postID),
		ChatID:    chatID,
		CreatedAt: createdAt,
	}
	_, err := arango.InsertDocument(ctx, r.db, arango.EdgeResponded, edge)
	return err
}

func (r *repository) CreateVotedEdge(ctx context.Context, userID, postID, option string, createdAt int64) error {
	edge := VotedEdge{
		From:      fmt.Sprintf("users/%s", userID),
		To:        fmt.Sprintf("posts/%s", postID),
		Option:    option,
		CreatedAt: createdAt,
	}
	_, err := arango.InsertDocument(ctx, r.db, arango.EdgeVoted, edge)
	return err
}

func (r *repository) CreatePostHasTagEdge(ctx context.Context, postID, tagKey string, confidence float64) error {
	edge := PostHasTagEdge{
		From:       fmt.Sprintf("posts/%s", postID),
		To:         fmt.Sprintf("tags/%s", tagKey),
		Confidence: confidence,
		Source:     "ai",
	}
	_, err := arango.InsertDocument(ctx, r.db, arango.EdgePostHasTag, edge)
	return err
}

func (r *repository) GetPostTags(ctx context.Context, postID string) ([]string, error) {
	return arango.Query[string](ctx, r.db, GetPostTags, map[string]any{
		"postId": fmt.Sprintf("posts/%s", postID),
	})
}

func (r *repository) GetVotes(ctx context.Context, postID string) (map[string]int, error) {
	type voteCount struct {
		Option string `json:"option"`
		Count  int    `json:"count"`
	}
	results, err := arango.Query[voteCount](ctx, r.db, GetPollVotes, map[string]any{
		"postId": fmt.Sprintf("posts/%s", postID),
	})
	if err != nil {
		return nil, err
	}

	votes := make(map[string]int)
	for _, v := range results {
		votes[v.Option] = v.Count
	}
	return votes, nil
}

func (r *repository) HasUserVoted(ctx context.Context, userID, postID string) (bool, error) {
	results, err := arango.Query[bool](ctx, r.db, CheckUserVoted, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
		"postId": fmt.Sprintf("posts/%s", postID),
	})
	if err != nil {
		return false, err
	}
	return len(results) > 0, nil
}

func (r *repository) HasUserResponded(ctx context.Context, userID, postID string) (bool, error) {
	results, err := arango.Query[bool](ctx, r.db, CheckUserResponded, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
		"postId": fmt.Sprintf("posts/%s", postID),
	})
	if err != nil {
		return false, err
	}
	return len(results) > 0, nil
}

func (r *repository) GetAuthor(ctx context.Context, postID string) (*PostAuthor, error) {
	return arango.QueryOne[PostAuthor](ctx, r.db, GetPostAuthor, map[string]any{
		"postId": fmt.Sprintf("posts/%s", postID),
	})
}
