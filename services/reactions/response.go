package reactions

import (
	"errors"

	"gorm.io/gorm"

	"thoughts_backend_api/models"
)

type reactionSummaryResponse struct {
	ThoughtID        uint    `json:"thought_id"`
	ThumbsUpCount    int64   `json:"thumbs_up_count"`
	ThumbsDownCount  int64   `json:"thumbs_down_count"`
	MyReaction       *string `json:"my_reaction,omitempty"`
}

func buildReactionSummary(db *gorm.DB, thoughtID uint, userID uint) (reactionSummaryResponse, error) {
	var thumbsUpCount int64
	if err := db.Model(&models.Reaction{}).
		Where("thought_id = ? AND type = ?", thoughtID, models.ReactionTypeThumbsUp).
		Count(&thumbsUpCount).Error; err != nil {
		return reactionSummaryResponse{}, err
	}

	var thumbsDownCount int64
	if err := db.Model(&models.Reaction{}).
		Where("thought_id = ? AND type = ?", thoughtID, models.ReactionTypeThumbsDown).
		Count(&thumbsDownCount).Error; err != nil {
		return reactionSummaryResponse{}, err
	}

	var myReaction *string
	var reaction models.Reaction
	if err := db.Where("thought_id = ? AND user_id = ?", thoughtID, userID).First(&reaction).Error; err == nil {
		reactionType := reaction.Type
		myReaction = &reactionType
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return reactionSummaryResponse{}, err
	}

	return reactionSummaryResponse{
		ThoughtID:       thoughtID,
		ThumbsUpCount:   thumbsUpCount,
		ThumbsDownCount: thumbsDownCount,
		MyReaction:      myReaction,
	}, nil
}
