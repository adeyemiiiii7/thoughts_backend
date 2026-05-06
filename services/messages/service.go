package messages

import (
	"thoughts_backend_api/models"
	"gorm.io/gorm"
)

func normalizeUserPair(a, b uint) (uint, uint) {
	if a < b {
		return a, b
	}
	return b, a
}

func getOrCreateConversation(db *gorm.DB, userAID, userBID uint) (*models.Conversation, error) {
	userOneID, userTwoID := normalizeUserPair(userAID, userBID)

	var conversation models.Conversation
	err := db.
		Where("user_one_id = ? AND user_two_id = ?", userOneID, userTwoID).
		First(&conversation).Error
	if err == nil {
		return &conversation, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	conversation = models.Conversation{
		UserOneID: userOneID,
		UserTwoID: userTwoID,
	}

	if err := db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func getConversationForUser(db *gorm.DB, conversationID, userID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := db.
		Where("id = ? AND (user_one_id = ? OR user_two_id = ?)", conversationID, userID, userID).
		First(&conversation).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func otherUserID(conversation models.Conversation, currentUserID uint) uint {
	if conversation.UserOneID == currentUserID {
		return conversation.UserTwoID
	}
	return conversation.UserOneID
}
