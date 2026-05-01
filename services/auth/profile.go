package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var currentUser models.User
	if err := h.db.Preload("Interests").First(&currentUser, user.ID).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load current user"})
		return
	}

	writeJSON(w, http.StatusOK, newProfileResponse(currentUser))
}

func (h *Handler) UpdateInterests(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req updateInterestsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	normalizedNames := normalizeInterestNames(req.Interests)
	interests := make([]models.Interest, 0, len(normalizedNames))

	for _, name := range normalizedNames {
		interest := models.Interest{Name: name}
		if err := h.db.FirstOrCreate(&interest, models.Interest{Name: name}).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save interests"})
			return
		}
		interests = append(interests, interest)
	}

	currentUser := models.User{ID: user.ID}
	if err := h.db.Model(&currentUser).Association("Interests").Replace(interests); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update interests"})
		return
	}

	if err := h.db.Preload("Interests").First(&currentUser, user.ID).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load updated user"})
		return
	}

	writeJSON(w, http.StatusOK, newProfileResponse(currentUser))
}

func newProfileResponse(user models.User) profileResponse {
	return profileResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Interests:     user.Interests,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func normalizeInterestNames(names []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(names))

	for _, name := range names {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}

		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	return result
}
