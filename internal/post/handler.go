package post

import (
	"net/http"

	"github.com/askme/api/pkg/httputil"
)

type handler struct {
	service Service
}

// NewHandler creates a new post handler
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// GetPost handles GET /posts/{postId}
func (h *handler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := httputil.PathValue(r, "postId")
	if postID == "" {
		httputil.Error(w, http.StatusBadRequest, "postId is required")
		return
	}

	post, err := h.service.GetPost(r.Context(), postID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, post)
}

// CreatePost handles POST /posts
func (h *handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	req, err := httputil.DecodeJSON[CreatePostRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AuthorID == "" {
		httputil.Error(w, http.StatusBadRequest, "authorId is required")
		return
	}
	if req.Text == "" {
		httputil.Error(w, http.StatusBadRequest, "text is required")
		return
	}

	resp, err := h.service.CreatePost(r.Context(), req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// CreatePoll handles POST /posts/poll
func (h *handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	req, err := httputil.DecodeJSON[CreatePostRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AuthorID == "" {
		httputil.Error(w, http.StatusBadRequest, "authorId is required")
		return
	}
	if req.Text == "" {
		httputil.Error(w, http.StatusBadRequest, "text is required")
		return
	}
	if len(req.PollOptions) < 2 {
		httputil.Error(w, http.StatusBadRequest, "at least 2 poll options are required")
		return
	}

	resp, err := h.service.CreatePoll(r.Context(), req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// RespondToPost handles POST /posts/{postId}/respond
func (h *handler) RespondToPost(w http.ResponseWriter, r *http.Request) {
	postID := httputil.PathValue(r, "postId")
	if postID == "" {
		httputil.Error(w, http.StatusBadRequest, "postId is required")
		return
	}

	req, err := httputil.DecodeJSON[RespondToPostRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}
	if req.Text == "" {
		httputil.Error(w, http.StatusBadRequest, "text is required")
		return
	}

	resp, err := h.service.RespondToPost(r.Context(), postID, req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// Vote handles POST /posts/{postId}/vote
func (h *handler) Vote(w http.ResponseWriter, r *http.Request) {
	postID := httputil.PathValue(r, "postId")
	if postID == "" {
		httputil.Error(w, http.StatusBadRequest, "postId is required")
		return
	}

	req, err := httputil.DecodeJSON[VoteRequest](r)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		httputil.Error(w, http.StatusBadRequest, "userId is required")
		return
	}
	if req.Option == "" {
		httputil.Error(w, http.StatusBadRequest, "option is required")
		return
	}

	resp, err := h.service.Vote(r.Context(), postID, req)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}
