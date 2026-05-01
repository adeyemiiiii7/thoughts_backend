package models

import "time"

const (
	ReactionTypeThumbsUp   = "thumbs_up"
	ReactionTypeThumbsDown = "thumbs_down"
)

type Reaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ThoughtID uint      `gorm:"not null;index;uniqueIndex:idx_user_thought_reaction" json:"thought_id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_user_thought_reaction" json:"user_id"`
	User      User      `json:"user,omitempty"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
