package messages

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	sender, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	recipientIDParam := chi.URLParam(r, "id")
	recipientID, err := strconv.ParseUint(recipientIDParam, 10, 64)
	if err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	if sender.ID == uint(recipientID) {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "you cannot message yourself"})
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Body = strings.TrimSpace(req.Body)
	if req.Body == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "body is required"})
		return
	}

	var recipient models.User
	if err := h.db.First(&recipient, uint(recipientID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "recipient not found"})
		return
	}

	allowed, err := canUsersMessage(h.db, sender.ID, recipient.ID)
	if err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to check messaging permission"})
		return
	}
	if !allowed {
		shared.RespondJSON(w, http.StatusForbidden, map[string]string{"error": "you can only message users who mutually follow you"})
		return
	}

	conversation, err := getOrCreateConversation(h.db, sender.ID, recipient.ID)
	if err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get conversation"})
		return
	}

	message := models.Message{
		ConversationID: conversation.ID,
		SenderID:       sender.ID,
		RecipientID:    recipient.ID,
		Body:           req.Body,
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&message).Error; err != nil {
			return err
		}
		return tx.Model(conversation).Updates(map[string]any{
			"last_message_at": message.CreatedAt,
			"updated_at":      time.Now(),
		}).Error
	}); err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save message"})
		return
	}

	payload := messageEvent{
		Type:           "message.created",
		ConversationID: conversation.ID,
		Message:        message,
	}
	_ = h.hub.SendToUser(r.Context(), recipient.ID, payload)
	_ = h.hub.SendToUser(r.Context(), sender.ID, payload)

	shared.RespondJSON(w, http.StatusCreated, payload)
}
