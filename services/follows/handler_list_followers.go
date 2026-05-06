package follows

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	pagination := shared.ParsePagination(r)

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

	var total int64
	if err := h.db.Model(&models.Follow{}).Where("following_id = ?", user.ID).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to count followers",
		})
		return
	}

	var follows []models.Follow
	if err := h.db.
		Where("following_id = ?", user.ID).
		Preload("Follower").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&follows).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load followers",
		})
		return
	}

	followers := make([]models.User, 0, len(follows))
	for _, follow := range follows {
		followers = append(followers, follow.Follower)
	}

	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(followers, pagination, total))
}
