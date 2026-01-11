package main

import (
	"github.com/askme/api/internal/chat"
	"github.com/askme/api/internal/feed"
	"github.com/askme/api/internal/post"
	"github.com/askme/api/internal/tag"
	"github.com/askme/api/internal/user"
	"github.com/askme/api/pkg/arango"
)

// App holds all feature module handlers (using interfaces for easy framework switching)
type App struct {
	userHandler user.Handler
	postHandler post.Handler
	chatHandler chat.Handler
	feedHandler feed.Handler
	tagHandler  tag.Handler
}

// NewApp initializes all feature modules with dependency injection
func NewApp(db *arango.Client) *App {
	// User feature
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Tag feature
	tagRepo := tag.NewRepository(db)
	tagService := tag.NewService(tagRepo)
	tagHandler := tag.NewHandler(tagService)

	// Post feature (depends on tag service for tag normalization)
	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo, tagService)
	postHandler := post.NewHandler(postService)

	// Chat feature
	chatRepo := chat.NewRepository(db)
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	// Feed feature (depends on post and chat repos for aggregation)
	feedRepo := feed.NewRepository(db)
	feedService := feed.NewService(feedRepo, postRepo, chatRepo)
	feedHandler := feed.NewHandler(feedService)

	return &App{
		userHandler: userHandler,
		postHandler: postHandler,
		chatHandler: chatHandler,
		feedHandler: feedHandler,
		tagHandler:  tagHandler,
	}
}
