package models

import "time"


type Conversation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserOneID     uint      `gorm:"not null;index" json:"user_one_id"`
	UserTwoID     uint      `gorm:"not null;index" json:"user_two_id"`
	Messages      []Message `json:"messages,omitempty"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
