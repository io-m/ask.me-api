package post

import "github.com/askme/api/internal/domain"

// Post represents a post document in ArangoDB
type Post struct {
	Key         string              `json:"_key,omitempty"`
	AuthorID    string              `json:"authorId"`
	PostType    domain.PostType     `json:"postType"`
	Text        string              `json:"text"`
	PollOptions []string            `json:"pollOptions,omitempty"`
	Category    domain.PostCategory `json:"category"`
	Intent      string              `json:"intent"`
	Depth       domain.PostDepth    `json:"depth"`
	AIRaw       domain.AIRawData    `json:"aiRaw,omitempty"`
	CreatedAt   int64               `json:"createdAt"`
}

// PostAuthor represents author information for API responses
type PostAuthor struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

// CreatedEdge represents the authorship edge
type CreatedEdge struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	CreatedAt int64  `json:"createdAt"`
}

// RespondedEdge represents a user's response to a post
type RespondedEdge struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	ChatID    string `json:"chatId,omitempty"`
	CreatedAt int64  `json:"createdAt"`
}

// VotedEdge represents a poll vote
type VotedEdge struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	Option    string `json:"option"`
	CreatedAt int64  `json:"createdAt"`
}

// PostHasTagEdge represents the post-tag relationship
type PostHasTagEdge struct {
	From       string  `json:"_from"`
	To         string  `json:"_to"`
	Confidence float64 `json:"confidence,omitempty"`
	Source     string  `json:"source,omitempty"`
}

// CreatePostRequest is the request payload for creating a post
type CreatePostRequest struct {
	AuthorID    string           `json:"authorId"`
	PostType    domain.PostType  `json:"postType"`
	Text        string           `json:"text"`
	PollOptions []string         `json:"pollOptions,omitempty"`
	AIRaw       domain.AIRawData `json:"aiRaw,omitempty"`
}

// CreatePostResponse is the response payload for creating a post
type CreatePostResponse struct {
	Key       string              `json:"_key"`
	Category  domain.PostCategory `json:"category"`
	Tags      []string            `json:"tags"`
	CreatedAt int64               `json:"createdAt"`
}

// RespondToPostRequest is the request payload for responding to a post
type RespondToPostRequest struct {
	UserID   string          `json:"userId"`
	Text     string          `json:"text"`
	ChatType domain.ChatType `json:"chatType,omitempty"`
}

// RespondToPostResponse is the response for responding to a post
type RespondToPostResponse struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"messageId"`
	CreatedAt int64  `json:"createdAt"`
}

// VoteRequest is the request payload for voting on a poll
type VoteRequest struct {
	UserID string `json:"userId"`
	Option string `json:"option"`
}

// VoteResponse is the response for voting on a poll
type VoteResponse struct {
	PostID string         `json:"postId"`
	Option string         `json:"option"`
	Votes  map[string]int `json:"votes"`
}

// GetPostResponse is the response for getting a post
type GetPostResponse struct {
	Key         string              `json:"_key"`
	AuthorID    string              `json:"authorId"`
	PostType    domain.PostType     `json:"postType"`
	Text        string              `json:"text"`
	PollOptions []string            `json:"pollOptions,omitempty"`
	Category    domain.PostCategory `json:"category"`
	Intent      string              `json:"intent"`
	Depth       domain.PostDepth    `json:"depth"`
	AIRaw       domain.AIRawData    `json:"aiRaw,omitempty"`
	Tags        []string            `json:"tags"`
	CreatedAt   int64               `json:"createdAt"`
}
