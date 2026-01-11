package chat

import "github.com/askme/api/internal/domain"

// Chat represents a chat document in ArangoDB
type Chat struct {
	Key              string          `json:"_key,omitempty"`
	PostID           string          `json:"postId"`
	Type             domain.ChatType `json:"type"`
	CreatedAt        int64           `json:"createdAt"`
	ExpiresAt        int64           `json:"expiresAt,omitempty"`
	ParticipantCount int             `json:"participantCount,omitempty"`
}

// Message represents a message document in ArangoDB
type Message struct {
	Key       string               `json:"_key,omitempty"`
	ChatID    string               `json:"chatId"`
	SenderID  string               `json:"senderId"`
	Text      string               `json:"text"`
	Status    domain.MessageStatus `json:"status"`
	CreatedAt int64                `json:"createdAt"`
}

// ParticipatesInEdge represents chat participation
type ParticipatesInEdge struct {
	From                 string                   `json:"_from"`
	To                   string                   `json:"_to"`
	Role                 domain.ParticipantRole   `json:"role"`
	Status               domain.ParticipantStatus `json:"status"`
	NotificationsEnabled bool                     `json:"notificationsEnabled"`
	JoinedAt             *int64                   `json:"joinedAt,omitempty"`
}

// TaggedEdge represents a user being tagged in a post
type TaggedEdge struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	CreatedAt int64  `json:"createdAt"`
}

// ChatPartner represents information about a chat partner
type ChatPartner struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

// QuestionContext represents the question that sparked a chat
type QuestionContext struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	AuthorID      string `json:"authorId"`
	CreatedAt     int64  `json:"createdAt"`
	FormattedTime string `json:"formattedTime"`
}

// LastMessage represents the most recent message in a chat
type LastMessage struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	SenderID      string `json:"senderId"`
	CreatedAt     int64  `json:"createdAt"`
	FormattedTime string `json:"formattedTime"`
}

// ChatThread represents a chat thread for listing
type ChatThread struct {
	ID          string          `json:"id"`
	Question    QuestionContext `json:"question"`
	Partner     ChatPartner     `json:"partner"`
	LastMessage LastMessage     `json:"lastMessage"`
	UnreadCount int             `json:"unreadCount"`
	HasUnread   bool            `json:"hasUnread"`
}

// ChatThreadsResponse is the response for listing chat threads
type ChatThreadsResponse struct {
	Threads    []ChatThread `json:"threads"`
	NextCursor *string      `json:"nextCursor,omitempty"`
}

// Participant represents a chat participant
type Participant struct {
	ID        string                   `json:"id"`
	Username  string                   `json:"username"`
	AvatarURL *string                  `json:"avatarUrl,omitempty"`
	Role      domain.ParticipantRole   `json:"role"`
	Status    domain.ParticipantStatus `json:"status"`
}

// GetChatResponse is the response for getting a chat
type GetChatResponse struct {
	Key       string            `json:"_key"`
	PostID    string            `json:"postId"`
	Type      domain.ChatType   `json:"type"`
	Messages  []MessageResponse `json:"messages"`
	CreatedAt int64             `json:"createdAt"`
}

// MessageResponse is a message in API responses
type MessageResponse struct {
	Key       string               `json:"_key"`
	SenderID  string               `json:"senderId"`
	Text      string               `json:"text"`
	Status    domain.MessageStatus `json:"status"`
	CreatedAt int64                `json:"createdAt"`
}

// ParticipantsResponse is the response for listing participants
type ParticipantsResponse struct {
	ChatID       string         `json:"chatId"`
	Type         domain.ChatType `json:"type"`
	Participants []Participant   `json:"participants"`
}

// SendMessageRequest is the request for sending a message
type SendMessageRequest struct {
	SenderID string `json:"senderId"`
	Text     string `json:"text"`
}

// SendMessageResponse is the response for sending a message
type SendMessageResponse struct {
	MessageID string `json:"messageId"`
	CreatedAt int64  `json:"createdAt"`
}

// AcceptChatRequest is the request for accepting a chat invite
type AcceptChatRequest struct {
	UserID string `json:"userId"`
}

// AcceptChatResponse is the response for accepting a chat invite
type AcceptChatResponse struct {
	Success bool                     `json:"success"`
	ChatID  string                   `json:"chatId"`
	Status  domain.ParticipantStatus `json:"status"`
}

// MuteChatRequest is the request for muting a chat
type MuteChatRequest struct {
	UserID string `json:"userId"`
}

// MuteChatResponse is the response for muting a chat
type MuteChatResponse struct {
	Success bool                     `json:"success"`
	ChatID  string                   `json:"chatId"`
	Status  domain.ParticipantStatus `json:"status"`
}
