package messages

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	pagination := shared.ParsePagination(r)

	conversationIDParam := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseUint(conversationIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
		return
	}

	conversation, err := getConversationForUser(h.db, uint(conversationID), user.ID)
	if err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "conversation not found"})
		return
	}

	var total int64
	if err := h.db.Model(&models.Message{}).Where("conversation_id = ?", conversation.ID).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count messages"})
		return
	}

	var messages []models.Message
	if err := h.db.
		Where("conversation_id = ?", conversation.ID).
		Order("created_at ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&messages).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load messages"})
		return
	}

	readAt := time.Now()
	_ = h.db.Model(&models.Message{}).
		Where("conversation_id = ? AND recipient_id = ? AND read_at IS NULL", conversation.ID, user.ID).
		Update("read_at", readAt).Error

	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(messages, pagination, total))
}
