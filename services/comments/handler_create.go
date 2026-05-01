package comments

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid thought id",
		})
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "content is required",
		})
		return
	}

	var thought models.Thought
	if err := h.db.First(&thought, uint(thoughtID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "thought not found",
		})
		return
	}

	// A normal comment belongs directly to a thought.
	comment := models.Comment{
		ThoughtID: thought.ID,
		UserID:    user.ID,
		Content:   req.Content,
	}

	if err := h.db.Create(&comment).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create comment",
		})
		return
	}

	if err := h.db.Preload("User").First(&comment, comment.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load created comment",
		})
		return
	}

	shared.RespondJSON(w, http.StatusCreated, comment)
}

// ReplyComment creates a reply under an existing comment.

func (h *Handler) ReplyComment(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	parentCommentIDParam := chi.URLParam(r, "id")
	parentCommentID, err := strconv.ParseUint(parentCommentIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid comment id",
		})
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "content is required",
		})
		return
	}

	var parentComment models.Comment
	if err := h.db.First(&parentComment, uint(parentCommentID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "parent comment not found",
		})
		return
	}

	parentCommentIDUint := uint(parentCommentID)
	// A reply is still a comment. The difference is that it points to a parent comment.
	comment := models.Comment{
		ThoughtID:        parentComment.ThoughtID,
		UserID:           user.ID,
		Content:          req.Content,
		ParentCommentID:  &parentCommentIDUint,
	}

	if err := h.db.Create(&comment).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create reply",
		})
		return
	}

	if err := h.db.Preload("User").First(&comment, comment.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load created reply",
		})
		return
	}

	shared.RespondJSON(w, http.StatusCreated, comment)
}
