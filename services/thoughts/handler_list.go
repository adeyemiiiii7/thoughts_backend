package thoughts

import (
	"net/http"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var thoughts []models.Thought
	// Preload works a bit like "include" in Prisma or Sequelize.
	// We are asking GORM to return each thought together with its related data.
	if err := h.db.
		Preload("User").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load thoughts"})
		return
	}

	shared.RespondJSON(w, http.StatusOK, thoughts)
}
