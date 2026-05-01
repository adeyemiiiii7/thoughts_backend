package auth

import (
	"time"

	"thoughts_backend_api/models"
)

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type updateInterestsRequest struct {
	Interests []string `json:"interests"`
}

type profileResponse struct {
	ID            uint              `json:"id"`
	Username      string            `json:"username"`
	Email         string            `json:"email"`
	EmailVerified bool              `json:"email_verified"`
	Interests     []models.Interest `json:"interests"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}
