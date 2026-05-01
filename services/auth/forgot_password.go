package auth

import (
	"encoding/json"
	"log"
	"net/http"
	netmail "net/mail"
	"strings"
	"time"

	"gorm.io/gorm"

	"thoughts_backend_api/models"
	emailservice "thoughts_backend_api/services/email"
)

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
		return
	}

	if _, err := netmail.ParseAddress(req.Email); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email address"})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "if an account exists for that email, a password reset link has been sent",
		})
		return
	}

	token := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     generateVerificationToken(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := h.db.Where("user_id = ? AND used = ?", user.ID, false).Delete(&models.PasswordResetToken{}).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to prepare password reset"})
		return
	}

	if err := h.db.Create(&token).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create password reset token"})
		return
	}

	log.Println("attempting to send password reset email to:", user.Email)
	if err := emailservice.SendPasswordResetEmail(user.Email, token.Token); err != nil {
		log.Println("failed to send password reset email:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to send password reset email"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "if an account exists for that email, a password reset link has been sent",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	if req.Token == "" || req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "token and new_password are required"})
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "new password must be at least 8 characters"})
		return
	}

	var token models.PasswordResetToken
	if err := h.db.Where("token = ? AND used = ?", req.Token, false).First(&token).Error; err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid or expired password reset token"})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password reset token has expired"})
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).
			Where("id = ?", token.UserID).
			Update("password", hashedPassword).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.PasswordResetToken{}).
			Where("id = ?", token.ID).
			Update("used", true).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to reset password"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "password reset successful"})
}
