package user

import (
	"context"
	"fmt"

	"github.com/askme/api/pkg/arango"
)

type repository struct {
	db *arango.Client
}

// NewRepository creates a new user repository
func NewRepository(db *arango.Client) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
	return arango.QueryOne[User](ctx, r.db, GetUserByID, map[string]any{"key": id})
}

func (r *repository) Create(ctx context.Context, user *User) (string, error) {
	return arango.InsertDocument(ctx, r.db, arango.CollectionUsers, user)
}

func (r *repository) Update(ctx context.Context, user *User) error {
	return arango.UpdateDocument(ctx, r.db, arango.CollectionUsers, user.Key, user)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return arango.DeleteDocument(ctx, r.db, arango.CollectionUsers, id)
}

func (r *repository) CreateFollow(ctx context.Context, followerID, followeeID string) (string, error) {
	edge := FollowsEdge{
		From: fmt.Sprintf("users/%s", followerID),
		To:   fmt.Sprintf("users/%s", followeeID),
	}
	return arango.InsertDocument(ctx, r.db, arango.EdgeFollows, edge)
}

func (r *repository) DeleteFollow(ctx context.Context, followerID, followeeID string) error {
	_, err := arango.Query[any](ctx, r.db, DeleteFollowEdge, map[string]any{
		"from": fmt.Sprintf("users/%s", followerID),
		"to":   fmt.Sprintf("users/%s", followeeID),
	})
	return err
}

func (r *repository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	results, err := arango.Query[bool](ctx, r.db, CheckIsFollowing, map[string]any{
		"from": fmt.Sprintf("users/%s", followerID),
		"to":   fmt.Sprintf("users/%s", followeeID),
	})
	if err != nil {
		return false, err
	}
	return len(results) > 0, nil
}

func (r *repository) AreMutualFollowers(ctx context.Context, userID1, userID2 string) (bool, error) {
	result, err := arango.QueryOne[bool](ctx, r.db, CheckMutualFollowers, map[string]any{
		"user1": fmt.Sprintf("users/%s", userID1),
		"user2": fmt.Sprintf("users/%s", userID2),
	})
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, nil
	}
	return *result, nil
}

func (r *repository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	result, err := arango.QueryOne[int](ctx, r.db, GetFollowerCount, map[string]any{
		"user": fmt.Sprintf("users/%s", userID),
	})
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return *result, nil
}

func (r *repository) GetFollowingCount(ctx context.Context, userID string) (int, error) {
	result, err := arango.QueryOne[int](ctx, r.db, GetFollowingCount, map[string]any{
		"user": fmt.Sprintf("users/%s", userID),
	})
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return *result, nil
}
