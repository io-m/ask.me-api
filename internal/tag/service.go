package tag

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/askme/api/internal/domain"
)

type service struct {
	repo Repository
}

// NewService creates a new tag service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTag(ctx context.Context, id string) (*Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get tag: %w", err)
	}
	if tag == nil {
		return nil, domain.ErrNotFound
	}
	return tag, nil
}

func (s *service) ListTags(ctx context.Context, limit, offset int) ([]Tag, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *service) SearchTags(ctx context.Context, query string, limit int) ([]Tag, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.Search(ctx, query, limit)
}

// NormalizeTags converts raw AI tags to canonical tag keys
// It creates new tags if they don't exist
func (s *service) NormalizeTags(ctx context.Context, rawTags []string) ([]string, error) {
	if len(rawTags) == 0 {
		return nil, nil
	}

	normalizedKeys := make([]string, len(rawTags))
	g, gCtx := errgroup.WithContext(ctx)

	for i, rawTag := range rawTags {
		i, rawTag := i, rawTag // capture loop variables
		g.Go(func() error {
			key, err := s.normalizeTag(gCtx, rawTag)
			if err != nil {
				return err
			}
			normalizedKeys[i] = key
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("normalize tags: %w", err)
	}

	// Remove duplicates
	seen := make(map[string]bool)
	unique := make([]string, 0, len(normalizedKeys))
	for _, key := range normalizedKeys {
		if key != "" && !seen[key] {
			seen[key] = true
			unique = append(unique, key)
		}
	}

	return unique, nil
}

func (s *service) normalizeTag(ctx context.Context, rawTag string) (string, error) {
	// First, try to find existing tag by alias
	existingTag, err := s.repo.GetByAlias(ctx, rawTag)
	if err != nil {
		return "", err
	}

	if existingTag != nil {
		// Increment usage count asynchronously (fire-and-forget)
		go func() {
			_ = s.repo.IncrementUsageCount(context.Background(), existingTag.Key)
		}()
		return existingTag.Key, nil
	}

	// Create new tag with normalized key
	key := toTagKey(rawTag)
	label := toTagLabel(rawTag)

	newTag := &Tag{
		Key:        key,
		Label:      label,
		Aliases:    []string{rawTag},
		UsageCount: 1,
		CreatedAt:  time.Now().UnixMilli(),
	}

	_, err = s.repo.Create(ctx, newTag)
	if err != nil {
		// If creation fails due to duplicate key, try to get existing
		existingTag, getErr := s.repo.GetByID(ctx, key)
		if getErr == nil && existingTag != nil {
			return existingTag.Key, nil
		}
		return "", err
	}

	return key, nil
}

// toTagKey converts a raw tag to a kebab-case key
// "Career Change" -> "career-change"
// "backend dev" -> "backend-dev"
func toTagKey(raw string) string {
	// Convert to lowercase
	s := strings.ToLower(raw)
	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]")
	s = reg.ReplaceAllString(s, "")
	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")
	// Trim hyphens from start and end
	s = strings.Trim(s, "-")
	return s
}

// toTagLabel converts a raw tag to a display label
// "career-change" -> "Career Change"
// "backend dev" -> "Backend Dev"
func toTagLabel(raw string) string {
	// Replace hyphens and underscores with spaces
	s := strings.ReplaceAll(raw, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	// Title case
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
