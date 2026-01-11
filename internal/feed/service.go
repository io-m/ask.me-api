package feed

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/askme/api/internal/chat"
	"github.com/askme/api/internal/domain"
	"github.com/askme/api/internal/post"
)

type service struct {
	repo     Repository
	postRepo post.Repository
	chatRepo chat.Repository
}

// NewService creates a new feed service
func NewService(repo Repository, postRepo post.Repository, chatRepo chat.Repository) Service {
	return &service{
		repo:     repo,
		postRepo: postRepo,
		chatRepo: chatRepo,
	}
}

func (s *service) GetFeed(ctx context.Context, query FeedQuery) (*FeedResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 50 {
		query.Limit = 50
	}

	// Validate optional filters
	if query.Category != "" && !domain.ValidCategory(query.Category) {
		query.Category = "" // Ignore invalid category
	}

	// Get user preferences in parallel for better personalization
	g, gCtx := errgroup.WithContext(ctx)

	var userTags []string
	var userCategories []string
	var userIntents []string

	g.Go(func() error {
		var err error
		userTags, err = s.repo.GetUserInteractionTags(gCtx, query.UserID)
		return err
	})

	g.Go(func() error {
		var err error
		userCategories, err = s.repo.GetUserCategories(gCtx, query.UserID)
		return err
	})

	g.Go(func() error {
		var err error
		userIntents, err = s.repo.GetUserIntents(gCtx, query.UserID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("get user preferences: %w", err)
	}

	// Log preferences for debugging (can be used for more sophisticated ranking)
	_ = userTags
	_ = userCategories
	_ = userIntents

	// Get recommended posts
	items, nextCursor, err := s.repo.GetRecommendedPosts(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get recommended posts: %w", err)
	}

	// Format timestamps
	for i := range items {
		if items[i].LastMessage != nil {
			items[i].LastMessage.FormattedTime = formatTime(items[i].LastMessage.CreatedAt)
		}
	}

	var cursorPtr *string
	if nextCursor != "" {
		cursorPtr = &nextCursor
	}

	return &FeedResponse{
		Items:      items,
		NextCursor: cursorPtr,
	}, nil
}

// formatTime formats a timestamp to a human-readable string
func formatTime(timestamp int64) string {
	t := time.UnixMilli(timestamp)
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 48*time.Hour:
		return "Yesterday"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	case diff < 14*24*time.Hour:
		return "1 week ago"
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / (24 * 7))
		return fmt.Sprintf("%d weeks ago", weeks)
	default:
		months := int(diff.Hours() / (24 * 30))
		if months <= 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}
