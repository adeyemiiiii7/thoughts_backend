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
	reactionIDParam := chi.URLParam(r, "id")
	reactionID, err := strconv.ParseUint(reactionIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid reaction id",
		})
		return
	}
	var reaction models.Reaction
	if err := h.db.First(&reaction, uint(reactionID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "reaction not found",
		})
		return
	}
	if reaction.UserID != user.ID {
		shared.RespondJSON(w, http.StatusForbidden, map[string]string{
			"error": "you can only delete your own reactions",
		})
		return
	}
	if err := h.db.Delete(&reaction).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete reaction",
		})
		return
	}
	shared.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "reaction deleted successfully",

	})
}