package chat

import (
	"context"
	"net/http"

	"github.com/askme/api/internal/domain"
)

// Repository defines the interface for chat data access
type Repository interface {
	GetByID(ctx context.Context, id string) (*Chat, error)
	Create(ctx context.Context, chat *Chat) (string, error)

	// Message operations
	CreateMessage(ctx context.Context, msg *Message) (string, error)
	GetMessages(ctx context.Context, chatID string, limit, offset int) ([]Message, error)
	UpdateMessageStatus(ctx context.Context, msgID string, status domain.MessageStatus) error
	GetUnreadCount(ctx context.Context, chatID, userID string) (int, error)

	// Participation operations
	CreateParticipation(ctx context.Context, edge *ParticipatesInEdge) error
	UpdateParticipation(ctx context.Context, userID, chatID string, status domain.ParticipantStatus, notificationsEnabled bool, joinedAt *int64) error
	GetParticipation(ctx context.Context, userID, chatID string) (*ParticipatesInEdge, error)
	GetParticipants(ctx context.Context, chatID string) ([]Participant, error)

	// Chat thread queries
	GetUserChatThreads(ctx context.Context, userID string, limit int, cursor string) ([]ChatThread, string, error)
	GetChatForPostAndUser(ctx context.Context, postID, userID string) (*Chat, error)

	// Reaction operations
	GetReaction(ctx context.Context, userID, messageID string) (*ReactedEdge, error)
	UpsertReaction(ctx context.Context, edge *ReactedEdge) error
	DeleteReaction(ctx context.Context, userID, messageID string) error
}

// Service defines the interface for chat business logic
type Service interface {
	GetChat(ctx context.Context, chatID string) (*GetChatResponse, error)
	GetUserChats(ctx context.Context, userID string, limit int, cursor string) (*ChatThreadsResponse, error)
	SendMessage(ctx context.Context, chatID string, req *SendMessageRequest) (*SendMessageResponse, error)
	AcceptChat(ctx context.Context, chatID string, req *AcceptChatRequest) (*AcceptChatResponse, error)
	MuteChat(ctx context.Context, chatID string, req *MuteChatRequest) (*MuteChatResponse, error)
	GetParticipants(ctx context.Context, chatID string) (*ParticipantsResponse, error)
	CreateChat(ctx context.Context, postID string, chatType domain.ChatType, participants []string) (string, error)
	ReactToMessage(ctx context.Context, req *ReactToMessageRequest) (*ReactToMessageResponse, error)
}

// Handler defines the interface for chat HTTP handlers
type Handler interface {
	GetUserChats(w http.ResponseWriter, r *http.Request)
	GetChat(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	AcceptChat(w http.ResponseWriter, r *http.Request)
	MuteChat(w http.ResponseWriter, r *http.Request)
	GetParticipants(w http.ResponseWriter, r *http.Request)
	ReactToMessage(w http.ResponseWriter, r *http.Request)
}
