package messages

import (
	"time"

	"thoughts_backend_api/models"
)

type sendMessageRequest struct {
	Body string `json:"body"`
}

type conversationResponse struct {
	ID            uint            `json:"id"`
	OtherUser     models.User     `json:"other_user"`
	LastMessage   *models.Message `json:"last_message,omitempty"`
	LastMessageAt time.Time       `json:"last_message_at"`
	UnreadCount   int64           `json:"unread_count"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type messageEvent struct {
	Type           string         `json:"type"`
	ConversationID uint           `json:"conversation_id"`
	Message        models.Message `json:"message"`
}
