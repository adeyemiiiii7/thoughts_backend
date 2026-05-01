package reactions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid thought id",
		})
		return
	}

	var req createReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	req.Type = strings.TrimSpace(strings.ToLower(req.Type))
	if req.Type != models.ReactionTypeThumbsUp && req.Type != models.ReactionTypeThumbsDown {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "type must be thumbs_up or thumbs_down",
		})
		return
	}

	var thought models.Thought
	if err := h.db.First(&thought, uint(thoughtID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "thought not found",
		})
		return
	}

	var reaction models.Reaction
	// One user should have only one reaction per thought.
	// If one already exists, update it instead of creating a duplicate.
	err = h.db.Where("thought_id = ? AND user_id = ?", thought.ID, user.ID).First(&reaction).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to load existing reaction",
			})
			return
		}

		reaction = models.Reaction{
			ThoughtID: thought.ID,
			UserID:    user.ID,
			Type:      req.Type,
		}

		if err := h.db.Create(&reaction).Error; err != nil {
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to create reaction",
			})
			return
		}
	} else {
		reaction.Type = req.Type
		if err := h.db.Save(&reaction).Error; err != nil {
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to update reaction",
			})
			return
		}
	}

	if err := h.db.Preload("User").First(&reaction, reaction.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load reaction",
		})
		return
	}

	shared.RespondJSON(w, http.StatusOK, reaction)
}
