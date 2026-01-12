package feed

import "github.com/askme/api/internal/domain"

// FeedAuthor represents author information in feed items
type FeedAuthor struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

// FeedLastMessage represents the last message preview in a feed item
type FeedLastMessage struct {
	ID            string               `json:"id"`
	Text          string               `json:"text"`
	SenderID      string               `json:"senderId"`
	Status        domain.MessageStatus `json:"status"`
	CreatedAt     int64                `json:"createdAt"`
	FormattedTime string               `json:"formattedTime"`
	MyReaction    *string              `json:"myReaction,omitempty"`
}

// FeedItem represents a single item in the feed
type FeedItem struct {
	ID          string              `json:"id"`
	PostType    domain.PostType     `json:"postType"`
	Text        string              `json:"text"`
	PollOptions []string            `json:"pollOptions,omitempty"`
	Category    domain.PostCategory `json:"category"`
	Intent      string              `json:"intent"`
	Depth       domain.PostDepth    `json:"depth"`
	Tags        []string            `json:"tags,omitempty"`
	Author      FeedAuthor          `json:"author"`
	ChatID      *string             `json:"chatId,omitempty"`
	LastMessage *FeedLastMessage    `json:"lastMessage,omitempty"`
	UnreadCount int                 `json:"unreadCount"`
	CreatedAt   int64               `json:"createdAt"`
}

// FeedResponse is the response for the feed endpoint
type FeedResponse struct {
	Items      []FeedItem `json:"items"`
	NextCursor *string    `json:"nextCursor,omitempty"`
}

// FeedQuery represents query parameters for the feed
type FeedQuery struct {
	UserID   string
	Limit    int
	Cursor   string
	Category string
	Depth    string
}
