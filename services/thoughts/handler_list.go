package thoughts

import (
	"net/http"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pagination := shared.ParsePagination(r)

	var total int64
	if err := h.db.Model(&models.Thought{}).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count thoughts"})
		return
	}

	var thoughts []models.Thought
	// Preload works a bit like "include" in Prisma or Sequelize.
	// We are asking GORM to return each thought together with its related data.
	if err := h.db.
		Preload("User").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load thoughts"})
		return
	}

	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(
		buildThoughtResponses(thoughts, nil),
		pagination,
		total,
	))
}
