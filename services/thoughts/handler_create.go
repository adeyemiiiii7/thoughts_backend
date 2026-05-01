package thoughts

import (
	"encoding/json"
	"net/http"
	"strings"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createThoughtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	if req.Content == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "content is required"})
		return
	}

	thought := models.Thought{
		UserID:  user.ID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.db.Create(&thought).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create thought"})
		return
	}

	if err := h.db.Preload("User").First(&thought, thought.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load created thought"})
		return
	}

	shared.RespondJSON(w, http.StatusCreated, thought)
}
