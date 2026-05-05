package follows

import (
	"errors"
	"net/http"
	"strconv"
 
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)
func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
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
	//prevent self-follow:
	if user.ID == uint(userID) {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cannot follow yourself",
		})
		return
	}
	// Check whether the follow already exists before creating a new row.
	var existingFollow models.Follow
	if err := h.db.Where("follower_id = ? AND following_id = ?", user.ID, uint(userID)).First(&existingFollow).Error; err == nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "already following this user",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to check follow relationship",
		})
		return
	}

	var targetUser models.User
	if err := h.db.First(&targetUser, uint(userID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
		return
	}

	follow := models.Follow{
		FollowerID:  user.ID,
		FollowingID: targetUser.ID,
	}

	if err := h.db.Create(&follow).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to follow user",
		})
		return
	}

	shared.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user followed successfully",
	})
}

