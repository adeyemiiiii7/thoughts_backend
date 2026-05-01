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

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	switch {
	case req.Username == "":
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username is required"})
		return
	case req.Email == "":
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
		return
	case req.Password == "":
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password is required"})
		return
	case len(req.Password) < 8:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password must be at least 8 characters"})
		return
	}

	if _, err := netmail.ParseAddress(req.Email); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email address"})
		return
	}

	var existingUser models.User
	if err := h.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email or username already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to check existing user"})
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      passwordHash,
		EmailVerified: false,
	}

	token := models.EmailVerificationToken{
		Token:     generateVerificationToken(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		token.UserID = user.ID
		if err := tx.Create(&token).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}

	log.Println("attempting to send verification email to:", user.Email)
	if err := emailservice.SendVerificationEmail(user.Email, token.Token); err != nil {
		log.Println("failed to send verification email:", err)
		writeJSON(w, http.StatusCreated, map[string]string{
			"message": "account created, but we could not send the verification email right now",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "signup successful, please check your email to verify your account",
	})
}
