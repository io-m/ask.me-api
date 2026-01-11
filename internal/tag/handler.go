package tag

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
)

type handler struct {
	service Service
}

// NewHandler creates a new tag handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetTag handles GET /tags/{tagId}
func (h *handler) GetTag(w http.ResponseWriter, r *http.Request) {
	tagID := httputil.PathValue(r, "tagId")
	if tagID == "" {
		httputil.Error(w, http.StatusBadRequest, "tagId is required")
		return
	}

	tag, err := h.service.GetTag(r.Context(), tagID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, tag)
}

// ListTags handles GET /tags
func (h *handler) ListTags(w http.ResponseWriter, r *http.Request) {
	limit := httputil.QueryInt(r, "limit", 50)
	offset := httputil.QueryInt(r, "offset", 0)
	query := httputil.QueryString(r, "q", "")

	var tags []Tag
	var err error

	if query != "" {
		tags, err = h.service.SearchTags(r.Context(), query, limit)
	} else {
		tags, err = h.service.ListTags(r.Context(), limit, offset)
	}

	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, tags)
}
