package feed

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
)

type handler struct {
	service Service
}

// NewHandler creates a new feed handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetFeed handles GET /feed
func (h *handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID := httputil.QueryString(r, "userId", "")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId query parameter is required")
		return
	}

	query := FeedQuery{
		UserID:   userID,
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
