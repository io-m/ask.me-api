package chat

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
)

type handler struct {
	service Service
}

// NewHandler creates a new chat handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetUserChats handles GET /users/{userId}/chats
func (h *handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID := httputil.PathValue(r, "userId")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	limit := httputil.QueryInt(r, "limit", 50)
	cursor := httputil.QueryString(r, "cursor", "")

	resp, err := h.service.GetUserChats(r.Context(), userID, limit, cursor)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

// GetChat handles GET /chats/{chatId}
func (h *handler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatID := httputil.PathValue(r, "chatId")
	if chatID == "" {
		httputil.Error(w, http.StatusBadRequest, "chatId is required")
		return
	}

	resp, err := h.service.GetChat(r.Context(), chatID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

// SendMessage handles POST /chats/{chatId}/message
func (h *handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatID := httputil.PathValue(r, "chatId")
	if chatID == "" {
		httputil.Error(w, http.StatusBadRequest, "chatId is required")
		return
	}

	req, err := httputil.DecodeJSON[SendMessageRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SenderID == "" {
		httputil.Error(w, http.StatusBadRequest, "senderId is required")
		return
	}
	if req.Text == "" {
		httputil.Error(w, http.StatusBadRequest, "text is required")
		return
	}

	resp, err := h.service.SendMessage(r.Context(), chatID, req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// AcceptChat handles POST /chats/{chatId}/accept
func (h *handler) AcceptChat(w http.ResponseWriter, r *http.Request) {
	chatID := httputil.PathValue(r, "chatId")
	if chatID == "" {
		httputil.Error(w, http.StatusBadRequest, "chatId is required")
		return
	}

	req, err := httputil.DecodeJSON[AcceptChatRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	resp, err := h.service.AcceptChat(r.Context(), chatID, req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

// MuteChat handles POST /chats/{chatId}/mute
func (h *handler) MuteChat(w http.ResponseWriter, r *http.Request) {
	chatID := httputil.PathValue(r, "chatId")
	if chatID == "" {
		httputil.Error(w, http.StatusBadRequest, "chatId is required")
		return
	}

	req, err := httputil.DecodeJSON[MuteChatRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	resp, err := h.service.MuteChat(r.Context(), chatID, req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

// GetParticipants handles GET /chats/{chatId}/participants
func (h *handler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	chatID := httputil.PathValue(r, "chatId")
	if chatID == "" {
		httputil.Error(w, http.StatusBadRequest, "chatId is required")
		return
	}

	resp, err := h.service.GetParticipants(r.Context(), chatID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}
