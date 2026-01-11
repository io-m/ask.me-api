package user

import (
	"context"
	"net/http"
)

// Repository defines the interface for user data access
type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (string, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error

	// Follow operations
	CreateFollow(ctx context.Context, followerID, followeeID string) (string, error)
	DeleteFollow(ctx context.Context, followerID, followeeID string) error
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	AreMutualFollowers(ctx context.Context, userID1, userID2 string) (bool, error)

	// Stats
	GetFollowerCount(ctx context.Context, userID string) (int, error)
	GetFollowingCount(ctx context.Context, userID string) (int, error)
}

// Service defines the interface for user business logic
type Service interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	FollowUser(ctx context.Context, followerID, followeeID string) (*FollowUserResponse, error)
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	AreMutualFollowers(ctx context.Context, userID1, userID2 string) (bool, error)
}

// Handler defines the interface for user HTTP handlers
type Handler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	FollowUser(w http.ResponseWriter, r *http.Request)
}
