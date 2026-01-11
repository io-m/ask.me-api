package user

import (
	"context"
	"fmt"
	"time"

	"github.com/askme/api/internal/domain"
)

type service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	now := time.Now().UnixMilli()

	user := &User{
		Username:  req.Username,
		CreatedAt: now,
		Interests: req.Interests,
		Settings:  req.Settings,
		Stats: UserStats{
			PostsCreated:   0,
			ResponsesGiven: 0,
		},
	}

	key, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &CreateUserResponse{
		Key:       key,
		CreatedAt: now,
	}, nil
}

func (s *service) FollowUser(ctx context.Context, followerID, followeeID string) (*FollowUserResponse, error) {
	// Check if already following
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return nil, fmt.Errorf("check following: %w", err)
	}
	if isFollowing {
		return nil, domain.ErrAlreadyExists
	}

	// Verify both users exist (can be parallelized with errgroup)
	_, err = s.repo.GetByID(ctx, followerID)
	if err != nil {
		return nil, fmt.Errorf("get follower: %w", err)
	}

	_, err = s.repo.GetByID(ctx, followeeID)
	if err != nil {
		return nil, fmt.Errorf("get followee: %w", err)
	}

	followID, err := s.repo.CreateFollow(ctx, followerID, followeeID)
	if err != nil {
		return nil, fmt.Errorf("create follow: %w", err)
	}

	return &FollowUserResponse{
		Success:  true,
		FollowID: followID,
	}, nil
}

func (s *service) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return s.repo.DeleteFollow(ctx, followerID, followeeID)
}

func (s *service) AreMutualFollowers(ctx context.Context, userID1, userID2 string) (bool, error) {
	return s.repo.AreMutualFollowers(ctx, userID1, userID2)
}
