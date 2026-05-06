package thoughts

import (
	"time"

	"thoughts_backend_api/models"
)

type thoughtResponse struct {
	ID              uint             `json:"id"`
	UserID          uint             `json:"user_id"`
	User            models.User      `json:"user,omitempty"`
	Title           string           `json:"title"`
	Content         string           `json:"content"`
	Comments        []models.Comment `json:"comments,omitempty"`
	ThumbsUpCount   int              `json:"thumbs_up_count"`
	ThumbsDownCount int              `json:"thumbs_down_count"`
	MyReaction      *string          `json:"my_reaction,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func buildThoughtResponses(thoughts []models.Thought, viewerID *uint) []thoughtResponse {
	responses := make([]thoughtResponse, 0, len(thoughts))
	for _, thought := range thoughts {
		responses = append(responses, buildThoughtResponse(thought, viewerID))
	}

	return responses
}

func buildThoughtResponse(thought models.Thought, viewerID *uint) thoughtResponse {
	thumbsUpCount := 0
	thumbsDownCount := 0
	var myReaction *string

	for _, reaction := range thought.Reactions {
		switch reaction.Type {
		case models.ReactionTypeThumbsUp:
			thumbsUpCount++
		case models.ReactionTypeThumbsDown:
			thumbsDownCount++
		}

		if viewerID != nil && reaction.UserID == *viewerID {
			reactionType := reaction.Type
			myReaction = &reactionType
		}
	}

	return thoughtResponse{
		ID:              thought.ID,
		UserID:          thought.UserID,
		User:            thought.User,
		Title:           thought.Title,
		Content:         thought.Content,
		Comments:        thought.Comments,
		ThumbsUpCount:   thumbsUpCount,
		ThumbsDownCount: thumbsDownCount,
		MyReaction:      myReaction,
		CreatedAt:       thought.CreatedAt,
		UpdatedAt:       thought.UpdatedAt,
	}
}
