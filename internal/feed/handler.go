package feed

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
	"github.com/askme/api/pkg/middleware"
)

type handler struct {
	service Service
}

// NewHandler creates a new feed handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetFeed handles GET /me/feed
func (h *handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	// Get current user from auth context
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := FeedQuery{
		UserID:   currentUserID,
		Limit:    httputil.QueryInt(r, "limit", 20),
		Cursor:   httputil.QueryString(r, "cursor", ""),
		Category: httputil.QueryString(r, "category", ""),
		Depth:    httputil.QueryString(r, "depth", ""),
	}

	resp, err := h.service.GetFeed(r.Context(), query)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}
