package comments

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListByThought(w http.ResponseWriter, r *http.Request) {
	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid thought id",
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

	var comments []models.Comment
	// This query means:
	// 1. Only load top-level comments for this thought.
	// 2. Also load the user who wrote each comment.
	// 3. Also load replies, and for each reply load its user too.
	if err := h.db.
		Where("thought_id = ? AND parent_comment_id IS NULL", thought.ID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load comments",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, comments)
}
