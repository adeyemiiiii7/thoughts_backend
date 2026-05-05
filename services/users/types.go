package users

import (
	"time"

	"thoughts_backend_api/models"
)

type publicProfileResponse struct {
	ID             uint              `json:"id"`
	Username       string            `json:"username"`
	Interests      []models.Interest `json:"interests"`
	ThoughtCount   int64             `json:"thought_count"`
	FollowerCount  int64             `json:"follower_count"`
	FollowingCount int64             `json:"following_count"`
	CreatedAt      time.Time         `json:"created_at"`
}
