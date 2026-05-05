package follows

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	userIDParam := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	var follow models.Follow
	if err := h.db.
		Where("follower_id = ? AND following_id = ?", user.ID, uint(userID)).
		First(&follow).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "follow relationship not found",
		})
		return
	}

	if err := h.db.Delete(&follow).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to unfollow user",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "user unfollowed successfully",
	})
}
