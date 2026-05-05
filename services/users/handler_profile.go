package users

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	var user models.User
	if err := h.db.Preload("Interests").First(&user, uint(userID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
		return
	}

	var thoughtCount int64
	if err := h.db.Model(&models.Thought{}).Where("user_id = ?", user.ID).Count(&thoughtCount).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load thought count",
		})
		return
	}

	var followerCount int64
	if err := h.db.Model(&models.Follow{}).Where("following_id = ?", user.ID).Count(&followerCount).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load follower count",
		})
		return
	}

	var followingCount int64
	if err := h.db.Model(&models.Follow{}).Where("follower_id = ?", user.ID).Count(&followingCount).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load following count",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, publicProfileResponse{
		ID:             user.ID,
		Username:       user.Username,
		Interests:      user.Interests,
		ThoughtCount:   thoughtCount,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		CreatedAt:      user.CreatedAt,
	})
}
