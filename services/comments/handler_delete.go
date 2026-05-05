package comments

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	commentIDParam := chi.URLParam(r, "id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid comment id",
		})
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, uint(commentID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "comment not found",
		})
		return
	}

	if comment.UserID != user.ID {
		shared.RespondJSON(w, http.StatusForbidden, map[string]string{
			"error": "you can only delete your own comments",
		})
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete comment",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "comment deleted successfully",
	})
}
