package messages

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	pagination := shared.ParsePagination(r)

	var total int64
	if err := h.db.Model(&models.Conversation{}).
		Where("user_one_id = ? OR user_two_id = ?", user.ID, user.ID).
		Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count conversations"})
		return
	}

	var conversations []models.Conversation
	if err := h.db.
		Where("user_one_id = ? OR user_two_id = ?", user.ID, user.ID).
		Order("last_message_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&conversations).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load conversations"})
		return
	}

	responses := make([]conversationResponse, 0, len(conversations))
	for _, conversation := range conversations {
		otherID := otherUserID(conversation, user.ID)

		var otherUser models.User
		if err := h.db.First(&otherUser, otherID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load conversation user"})
			return
		}

		var lastMessage models.Message
		var lastMessagePtr *models.Message
		if err := h.db.
			Where("conversation_id = ?", conversation.ID).
			Order("created_at DESC").
			First(&lastMessage).Error; err == nil {
			lastMessagePtr = &lastMessage
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load last message"})
			return
		}

		var unreadCount int64
		if err := h.db.Model(&models.Message{}).
			Where("conversation_id = ? AND recipient_id = ? AND read_at IS NULL", conversation.ID, user.ID).
			Count(&unreadCount).Error; err != nil {
			shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count unread messages"})
			return
		}

		responses = append(responses, conversationResponse{
			ID:            conversation.ID,
			OtherUser:     otherUser,
			LastMessage:   lastMessagePtr,
			LastMessageAt: conversation.LastMessageAt,
			UnreadCount:   unreadCount,
			CreatedAt:     conversation.CreatedAt,
			UpdatedAt:     conversation.UpdatedAt,
		})
	}

	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(responses, pagination, total))
}
