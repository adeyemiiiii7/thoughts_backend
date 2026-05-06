package reactions

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

	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid thought id",
		})
		return
	}

	var reaction models.Reaction
	if err := h.db.
		Where("thought_id = ? AND user_id = ?", uint(thoughtID), user.ID).
		First(&reaction).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "reaction not found",
		})
		return
	}

	if err := h.db.Delete(&reaction).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete reaction",
		})
		return
	}

	summary, err := buildReactionSummary(h.db, uint(thoughtID), user.ID)
	if err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load reaction summary",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, summary)
}
