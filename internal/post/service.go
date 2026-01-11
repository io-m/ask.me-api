package post

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/askme/api/internal/domain"
	"github.com/askme/api/internal/tag"
)

type service struct {
	repo       Repository
	tagService tag.Service
}

// NewService creates a new post service
func NewService(repo Repository, tagService tag.Service) Service {
	return &service{
		repo:       repo,
		tagService: tagService,
	}
}

func (s *service) GetPost(ctx context.Context, id string) (*GetPostResponse, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	if post == nil {
		return nil, domain.ErrNotFound
	}

	tags, err := s.repo.GetPostTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post tags: %w", err)
	}

	return &GetPostResponse{
		Key:         post.Key,
		AuthorID:    post.AuthorID,
		PostType:    post.PostType,
		Text:        post.Text,
		PollOptions: post.PollOptions,
		Category:    post.Category,
		Intent:      post.Intent,
		Depth:       post.Depth,
		AIRaw:       post.AIRaw,
		Tags:        tags,
		CreatedAt:   post.CreatedAt,
	}, nil
}

func (s *service) CreatePost(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error) {
	return s.createPostInternal(ctx, req, domain.PostTypeText)
}

func (s *service) CreatePoll(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error) {
	if len(req.PollOptions) < 2 {
		return nil, fmt.Errorf("%w: poll requires at least 2 options", domain.ErrInvalidInput)
	}
	return s.createPostInternal(ctx, req, domain.PostTypePoll)
}

func (s *service) createPostInternal(ctx context.Context, req *CreatePostRequest, postType domain.PostType) (*CreatePostResponse, error) {
	now := time.Now().UnixMilli()

	// Normalize AI-provided values
	category := domain.NormalizeCategory(req.AIRaw.Category)
	depth := domain.NormalizeDepth(req.AIRaw.Depth)

	post := &Post{
		AuthorID:    req.AuthorID,
		PostType:    postType,
		Text:        req.Text,
		PollOptions: req.PollOptions,
		Category:    category,
		Intent:      req.AIRaw.Intent,
		Depth:       depth,
		AIRaw:       req.AIRaw,
		CreatedAt:   now,
	}

	postKey, err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	// Create edges and normalize tags in parallel
	g, gCtx := errgroup.WithContext(ctx)

	// Create authorship edge
	g.Go(func() error {
		return s.repo.CreateCreatedEdge(gCtx, req.AuthorID, postKey, now)
	})

	// Normalize and create tag edges
	var normalizedTags []string
	g.Go(func() error {
		var tagErr error
		normalizedTags, tagErr = s.tagService.NormalizeTags(gCtx, req.AIRaw.Tags)
		if tagErr != nil {
			return tagErr
		}

		// Create tag edges (can be further parallelized if needed)
		for _, tagKey := range normalizedTags {
			if err := s.repo.CreatePostHasTagEdge(gCtx, postKey, tagKey, req.AIRaw.Confidence); err != nil {
				return err
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("post creation edges: %w", err)
	}

	return &CreatePostResponse{
		Key:       postKey,
		Category:  category,
		Tags:      normalizedTags,
		CreatedAt: now,
	}, nil
}

func (s *service) RespondToPost(ctx context.Context, postID string, req *RespondToPostRequest) (*RespondToPostResponse, error) {
	// Check if post exists
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	if post == nil {
		return nil, domain.ErrNotFound
	}

	// Check if user has already responded
	hasResponded, err := s.repo.HasUserResponded(ctx, req.UserID, postID)
	if err != nil {
		return nil, fmt.Errorf("check responded: %w", err)
	}
	if hasResponded {
		return nil, domain.ErrAlreadyExists
	}

	now := time.Now().UnixMilli()

	// TODO: Create or get chat, create message
	// This would involve the chat service
	chatID := ""    // Placeholder
	messageID := "" // Placeholder

	// Create responded edge
	if err := s.repo.CreateRespondedEdge(ctx, req.UserID, postID, chatID, now); err != nil {
		return nil, fmt.Errorf("create responded edge: %w", err)
	}

	return &RespondToPostResponse{
		ChatID:    chatID,
		MessageID: messageID,
		CreatedAt: now,
	}, nil
}

func (s *service) Vote(ctx context.Context, postID string, req *VoteRequest) (*VoteResponse, error) {
	// Check if post exists and is a poll
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	if post == nil {
		return nil, domain.ErrNotFound
	}
	if post.PostType != domain.PostTypePoll {
		return nil, fmt.Errorf("%w: post is not a poll", domain.ErrInvalidInput)
	}

	// Validate option
	validOption := false
	for _, opt := range post.PollOptions {
		if opt == req.Option {
			validOption = true
			break
		}
	}
	if !validOption {
		return nil, fmt.Errorf("%w: invalid poll option", domain.ErrInvalidInput)
	}

	// Check if user has already voted
	hasVoted, err := s.repo.HasUserVoted(ctx, req.UserID, postID)
	if err != nil {
		return nil, fmt.Errorf("check voted: %w", err)
	}
	if hasVoted {
		return nil, domain.ErrAlreadyExists
	}

	now := time.Now().UnixMilli()

	// Create vote edge
	if err := s.repo.CreateVotedEdge(ctx, req.UserID, postID, req.Option, now); err != nil {
		return nil, fmt.Errorf("create vote: %w", err)
	}

	// Get updated vote counts
	votes, err := s.repo.GetVotes(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get votes: %w", err)
	}

	return &VoteResponse{
		PostID: postID,
		Option: req.Option,
		Votes:  votes,
	}, nil
}
