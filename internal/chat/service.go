package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/askme/api/internal/domain"
)

type service struct {
	repo Repository
}

// NewService creates a new chat service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetChat(ctx context.Context, chatID string) (*GetChatResponse, error) {
	chat, err := s.repo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat: %w", err)
	}
	if chat == nil {
		return nil, domain.ErrNotFound
	}

	messages, err := s.repo.GetMessages(ctx, chatID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	msgResponses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		msgResponses[i] = MessageResponse{
			Key:       msg.Key,
			SenderID:  msg.SenderID,
			Text:      msg.Text,
			Status:    msg.Status,
			CreatedAt: msg.CreatedAt,
		}
	}

	return &GetChatResponse{
		Key:       chat.Key,
		PostID:    chat.PostID,
		Type:      chat.Type,
		Messages:  msgResponses,
		CreatedAt: chat.CreatedAt,
	}, nil
}

func (s *service) GetUserChats(ctx context.Context, userID string, limit int, cursor string) (*ChatThreadsResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	threads, nextCursor, err := s.repo.GetUserChatThreads(ctx, userID, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("get user chats: %w", err)
	}

	// Format times
	for i := range threads {
		threads[i].Question.FormattedTime = formatTime(threads[i].Question.CreatedAt)
		threads[i].LastMessage.FormattedTime = formatTime(threads[i].LastMessage.CreatedAt)
	}

	var cursorPtr *string
	if nextCursor != "" {
		cursorPtr = &nextCursor
	}

	return &ChatThreadsResponse{
		Threads:    threads,
		NextCursor: cursorPtr,
	}, nil
}

func (s *service) SendMessage(ctx context.Context, chatID string, req *SendMessageRequest) (*SendMessageResponse, error) {
	// Verify chat exists
	chat, err := s.repo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat: %w", err)
	}
	if chat == nil {
		return nil, domain.ErrNotFound
	}

	// Verify user is a participant
	participation, err := s.repo.GetParticipation(ctx, req.SenderID, chatID)
	if err != nil {
		return nil, fmt.Errorf("get participation: %w", err)
	}
	if participation == nil {
		return nil, domain.ErrForbidden
	}

	now := time.Now().UnixMilli()

	msg := &Message{
		ChatID:    fmt.Sprintf("chats/%s", chatID),
		SenderID:  fmt.Sprintf("users/%s", req.SenderID),
		Text:      req.Text,
		Status:    domain.MessageStatusSent,
		CreatedAt: now,
	}

	msgID, err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return &SendMessageResponse{
		MessageID: msgID,
		CreatedAt: now,
	}, nil
}

func (s *service) AcceptChat(ctx context.Context, chatID string, req *AcceptChatRequest) (*AcceptChatResponse, error) {
	// Verify participation exists and is pending
	participation, err := s.repo.GetParticipation(ctx, req.UserID, chatID)
	if err != nil {
		return nil, fmt.Errorf("get participation: %w", err)
	}
	if participation == nil {
		return nil, domain.ErrNotFound
	}
	if participation.Status != domain.StatusPending {
		return nil, fmt.Errorf("%w: chat already accepted or muted", domain.ErrInvalidInput)
	}

	now := time.Now().UnixMilli()
	err = s.repo.UpdateParticipation(ctx, req.UserID, chatID, domain.StatusActive, true, &now)
	if err != nil {
		return nil, fmt.Errorf("update participation: %w", err)
	}

	return &AcceptChatResponse{
		Success: true,
		ChatID:  chatID,
		Status:  domain.StatusActive,
	}, nil
}

func (s *service) MuteChat(ctx context.Context, chatID string, req *MuteChatRequest) (*MuteChatResponse, error) {
	// Verify participation exists
	participation, err := s.repo.GetParticipation(ctx, req.UserID, chatID)
	if err != nil {
		return nil, fmt.Errorf("get participation: %w", err)
	}
	if participation == nil {
		return nil, domain.ErrNotFound
	}

	err = s.repo.UpdateParticipation(ctx, req.UserID, chatID, domain.StatusMuted, false, participation.JoinedAt)
	if err != nil {
		return nil, fmt.Errorf("update participation: %w", err)
	}

	return &MuteChatResponse{
		Success: true,
		ChatID:  chatID,
		Status:  domain.StatusMuted,
	}, nil
}

func (s *service) GetParticipants(ctx context.Context, chatID string) (*ParticipantsResponse, error) {
	chat, err := s.repo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat: %w", err)
	}
	if chat == nil {
		return nil, domain.ErrNotFound
	}

	participants, err := s.repo.GetParticipants(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get participants: %w", err)
	}

	return &ParticipantsResponse{
		ChatID:       chatID,
		Type:         chat.Type,
		Participants: participants,
	}, nil
}

func (s *service) CreateChat(ctx context.Context, postID string, chatType domain.ChatType, participants []string) (string, error) {
	now := time.Now().UnixMilli()

	chat := &Chat{
		PostID:           fmt.Sprintf("posts/%s", postID),
		Type:             chatType,
		CreatedAt:        now,
		ParticipantCount: len(participants),
	}

	chatID, err := s.repo.Create(ctx, chat)
	if err != nil {
		return "", fmt.Errorf("create chat: %w", err)
	}

	// Create participation edges for all participants
	for i, userID := range participants {
		role := domain.RoleResponder
		status := domain.StatusActive
		if i == 0 {
			role = domain.RoleAuthor
		}

		edge := &ParticipatesInEdge{
			From:                 fmt.Sprintf("users/%s", userID),
			To:                   fmt.Sprintf("chats/%s", chatID),
			Role:                 role,
			Status:               status,
			NotificationsEnabled: true,
			JoinedAt:             &now,
		}

		if err := s.repo.CreateParticipation(ctx, edge); err != nil {
			return "", fmt.Errorf("create participation: %w", err)
		}
	}

	return chatID, nil
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
