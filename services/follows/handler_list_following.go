package follows

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, uint(userID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
		return
	}

	var follows []models.Follow
	if err := h.db.
		Where("follower_id = ?", user.ID).
		Preload("Following").
		Order("created_at DESC").
		Find(&follows).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load following users",
		})
		return
	}

	followingUsers := make([]models.User, 0, len(follows))
	for _, follow := range follows {
		followingUsers = append(followingUsers, follow.Following)
	}

	shared.RespondJSON(w, http.StatusOK, followingUsers)
}
