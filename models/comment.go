package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ThoughtID uint      `gorm:"not null;index" json:"thought_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `json:"user,omitempty"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
