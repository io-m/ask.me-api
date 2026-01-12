package user

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
	"github.com/askme/api/pkg/middleware"
)

type handler struct {
	service Service
}

// NewHandler creates a new user handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetUser handles GET /users/{userId}
func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := httputil.PathValue(r, "userId")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req, err := httputil.DecodeJSON[CreateUserRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" {
		httputil.Error(w, http.StatusBadRequest, "username is required")
		return
	}

	resp, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// FollowUser handles POST /me/follow/{userId}
func (h *handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from auth context
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user to follow from path
	followUserID := httputil.PathValue(r, "userId")
	if followUserID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	resp, err := h.service.FollowUser(r.Context(), currentUserID, followUserID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}
