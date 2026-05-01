package thoughts

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
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid thought id"})
		return
	}

	var thought models.Thought
	if err := h.db.First(&thought, uint(thoughtID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "thought not found"})
		return
	}

	if thought.UserID != user.ID {
		shared.RespondJSON(w, http.StatusForbidden, map[string]string{"error": "you can only delete your own thoughts"})
		return
	}

	if err := h.db.Delete(&thought).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete thought"})
		return
	}

	shared.RespondJSON(w, http.StatusOK, map[string]string{"message": "thought deleted successfully"})
}
