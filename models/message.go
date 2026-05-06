package models

import "time"

type Message struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ConversationID uint       `gorm:"not null;index" json:"conversation_id"`
	SenderID       uint       `gorm:"not null;index" json:"sender_id"`
	RecipientID    uint       `gorm:"not null;index" json:"recipient_id"`
	Body           string     `gorm:"not null" json:"body"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
