package thoughts

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
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
	if err := h.db.Model(&models.Thought{}).Where("user_id = ?", user.ID).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to count user thoughts",
		})
		return
	}

	var thoughts []models.Thought
	// This is like filtering posts by author in Prisma:
	// where: { userId: user.id }
	if err := h.db.
		Where("user_id = ?", user.ID).
		Preload("User").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load user thoughts",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(
		buildThoughtResponses(thoughts, nil),
		pagination,
		total,
	))
}
