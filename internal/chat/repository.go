package chat

import (
	"context"
	"fmt"

	"github.com/askme/api/internal/domain"
	"github.com/askme/api/pkg/arango"
)

type repository struct {
	db *arango.Client
}

// NewRepository creates a new chat repository
func NewRepository(db *arango.Client) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*Chat, error) {
	return arango.QueryOne[Chat](ctx, r.db, GetChatByID, map[string]any{"key": id})
}

func (r *repository) Create(ctx context.Context, chat *Chat) (string, error) {
	return arango.InsertDocument(ctx, r.db, arango.CollectionChats, chat)
}

func (r *repository) CreateMessage(ctx context.Context, msg *Message) (string, error) {
	return arango.InsertDocument(ctx, r.db, arango.CollectionMessages, msg)
}

func (r *repository) GetMessages(ctx context.Context, chatID string, limit, offset int) ([]Message, error) {
	return arango.Query[Message](ctx, r.db, GetChatMessages, map[string]any{
		"chatId": fmt.Sprintf("chats/%s", chatID),
		"limit":  limit,
		"offset": offset,
	})
}

func (r *repository) UpdateMessageStatus(ctx context.Context, msgID string, status domain.MessageStatus) error {
	_, err := arango.Query[any](ctx, r.db, UpdateMessageStatus, map[string]any{
		"key":    msgID,
		"status": status,
	})
	return err
}

func (r *repository) GetUnreadCount(ctx context.Context, chatID, userID string) (int, error) {
	result, err := arango.QueryOne[int](ctx, r.db, GetChatUnreadCount, map[string]any{
		"chatId": fmt.Sprintf("chats/%s", chatID),
		"userId": fmt.Sprintf("users/%s", userID),
	})
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return *result, nil
}

func (r *repository) CreateParticipation(ctx context.Context, edge *ParticipatesInEdge) error {
	_, err := arango.InsertDocument(ctx, r.db, arango.EdgeParticipatesIn, edge)
	return err
}

func (r *repository) UpdateParticipation(ctx context.Context, userID, chatID string, status domain.ParticipantStatus, notificationsEnabled bool, joinedAt *int64) error {
	_, err := arango.Query[any](ctx, r.db, UpdateParticipation, map[string]any{
		"from":                 fmt.Sprintf("users/%s", userID),
		"to":                   fmt.Sprintf("chats/%s", chatID),
		"status":               status,
		"notificationsEnabled": notificationsEnabled,
		"joinedAt":             joinedAt,
	})
	return err
}

func (r *repository) GetParticipation(ctx context.Context, userID, chatID string) (*ParticipatesInEdge, error) {
	return arango.QueryOne[ParticipatesInEdge](ctx, r.db, GetParticipation, map[string]any{
		"from": fmt.Sprintf("users/%s", userID),
		"to":   fmt.Sprintf("chats/%s", chatID),
	})
}

func (r *repository) GetParticipants(ctx context.Context, chatID string) ([]Participant, error) {
	return arango.Query[Participant](ctx, r.db, GetChatParticipants, map[string]any{
		"chatId": fmt.Sprintf("chats/%s", chatID),
	})
}

func (r *repository) GetUserChatThreads(ctx context.Context, userID string, limit int, cursor string) ([]ChatThread, string, error) {
	threads, err := arango.Query[ChatThread](ctx, r.db, GetUserChatThreads, map[string]any{
		"userId": fmt.Sprintf("users/%s", userID),
		"limit":  limit,
	})
	if err != nil {
		return nil, "", err
	}

	// Generate next cursor if there are more results
	var nextCursor string
	if len(threads) == limit && len(threads) > 0 {
		// Use last message createdAt as cursor
		nextCursor = fmt.Sprintf("%d", threads[len(threads)-1].LastMessage.CreatedAt)
	}

	return threads, nextCursor, nil
}

func (r *repository) GetChatForPostAndUser(ctx context.Context, postID, userID string) (*Chat, error) {
	return arango.QueryOne[Chat](ctx, r.db, GetChatForPostAndUser, map[string]any{
		"postId": fmt.Sprintf("posts/%s", postID),
		"userId": fmt.Sprintf("users/%s", userID),
	})
}
