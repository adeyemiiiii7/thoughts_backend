package auth

import (
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"thoughts_backend_api/models"
)

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenValue := strings.TrimSpace(r.URL.Query().Get("token"))
	if tokenValue == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "token is required"})
		return
	}

	var token models.EmailVerificationToken
	if err := h.db.Where("token = ? AND used = ?", tokenValue, false).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "invalid verification token"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load verification token"})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "verification token has expired"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).
			Where("id = ?", token.UserID).
			Update("email_verified", true).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.EmailVerificationToken{}).
			Where("id = ?", token.ID).
			Update("used", true).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to verify email"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "email verified successfully"})
}
