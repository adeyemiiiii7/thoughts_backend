package thoughts

import (
	"math/rand"
	"net/http"
	"sort"
	"time"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"
)

type feedCandidate struct {
	thought models.Thought
	score   float64
}

func (h *Handler) FollowingFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	pagination := shared.ParsePagination(r)

	var followedIDs []uint
	if err := h.db.Model(&models.Follow{}).
		Where("follower_id = ?", user.ID).
		Pluck("following_id", &followedIDs).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load follow graph",
		})
		return
	}

	authorIDs := append([]uint{user.ID}, followedIDs...)

	var total int64
	if err := h.db.Model(&models.Thought{}).Where("user_id IN ?", authorIDs).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to count following feed",
		})
		return
	}

	candidateLimit := pagination.Limit * 3
	candidateOffset := pagination.Offset * 3

	var thoughts []models.Thought
	if err := h.db.
		Where("user_id IN ?", authorIDs).
		Preload("User").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
		Offset(candidateOffset).
		Limit(candidateLimit).
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load following feed",
		})
		return
	}

	candidates := make([]feedCandidate, 0, len(thoughts))
	for _, thought := range thoughts {
		score := recencyScore(thought.CreatedAt) + float64(len(thought.Comments)*2+len(thought.Reactions))
		if thought.UserID == user.ID {
			score += 40
		}

		candidates = append(candidates, feedCandidate{thought: thought, score: score})
	}

	feed := diversifyFeed(candidates, pagination.Limit)
	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(
		buildThoughtResponses(feed, &user.ID),
		pagination,
		total,
	))
}

func (h *Handler) FYPFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	pagination := shared.ParsePagination(r)
	candidateLimit := pagination.Limit * 5
	if candidateLimit < 60 {
		candidateLimit = 60
	}
	if candidateLimit > 180 {
		candidateLimit = 180
	}
	candidateOffset := pagination.Offset * 5

	var viewer models.User
	if err := h.db.Preload("Interests").First(&viewer, user.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load current user",
		})
		return
	}

	var followedIDs []uint
	if err := h.db.Model(&models.Follow{}).
		Where("follower_id = ?", viewer.ID).
		Pluck("following_id", &followedIDs).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load follow graph",
		})
		return
	}

	var total int64
	if err := h.db.Model(&models.Thought{}).Where("user_id <> ?", viewer.ID).Count(&total).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to count fyp feed",
		})
		return
	}

	var thoughts []models.Thought
	if err := h.db.
		Preload("User").
		Preload("User.Interests").
		Preload("Comments").
		Preload("Reactions").
		Where("user_id <> ?", viewer.ID).
		Order("created_at DESC").
		Offset(candidateOffset).
		Limit(candidateLimit).
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load fyp candidates",
		})
		return
	}

	viewerInterestSet := make(map[uint]struct{}, len(viewer.Interests))
	for _, interest := range viewer.Interests {
		viewerInterestSet[interest.ID] = struct{}{}
	}

	followedSet := make(map[uint]struct{}, len(followedIDs))
	for _, followedID := range followedIDs {
		followedSet[followedID] = struct{}{}
	}

	// FYP is intentionally more exploratory than Following.
	// Shared interests are the strongest personalization signal here,
	// while recency, engagement, and a small randomness factor keep the page fresh.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	candidates := make([]feedCandidate, 0, len(thoughts))
	for _, thought := range thoughts {
		sharedInterests := countSharedInterests(viewerInterestSet, thought.User.Interests)
		engagementScore := len(thought.Comments)*2 + len(thought.Reactions)
		if engagementScore > 20 {
			engagementScore = 20
		}

		score := 0.0
		score += recencyScore(thought.CreatedAt)
		score += float64(sharedInterests) * 140
		score += float64(engagementScore) * 3
		score += rng.Float64() * 35

		if _, isFollowed := followedSet[thought.UserID]; isFollowed {
			score -= 120
		}

		if sharedInterests == 0 {
			score += 15
		}

		candidates = append(candidates, feedCandidate{thought: thought, score: score})
	}

	feed := diversifyFeed(candidates, pagination.Limit)
	shared.RespondJSON(w, http.StatusOK, shared.NewPaginatedResponse(
		buildThoughtResponses(feed, &user.ID),
		pagination,
		total,
	))
}

func countSharedInterests(viewerInterestSet map[uint]struct{}, interests []models.Interest) int {
	count := 0
	for _, interest := range interests {
		if _, ok := viewerInterestSet[interest.ID]; ok {
			count++
		}
	}

	return count
}

func recencyScore(createdAt time.Time) float64 {
	hoursOld := time.Since(createdAt).Hours()

	switch {
	case hoursOld <= 6:
		return 60
	case hoursOld <= 24:
		return 40
	case hoursOld <= 72:
		return 20
	default:
		return 0
	}
}

func diversifyFeed(candidates []feedCandidate, limit int) []models.Thought {
	if len(candidates) == 0 {
		return []models.Thought{}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].thought.CreatedAt.After(candidates[j].thought.CreatedAt)
		}
		return candidates[i].score > candidates[j].score
	})

	selected := make([]models.Thought, 0, limit)
	authorPickCount := make(map[uint]int)
	used := make([]bool, len(candidates))

	for len(selected) < limit {
		bestIndex := -1
		bestScore := -1.0

		for i, candidate := range candidates {
			if used[i] {
				continue
			}

			adjustedScore := candidate.score - float64(authorPickCount[candidate.thought.UserID]*75)
			if bestIndex == -1 ||
				adjustedScore > bestScore ||
				(adjustedScore == bestScore && candidate.thought.CreatedAt.After(candidates[bestIndex].thought.CreatedAt)) {
				bestIndex = i
				bestScore = adjustedScore
			}
		}

		if bestIndex == -1 {
			break
		}

		used[bestIndex] = true
		selected = append(selected, candidates[bestIndex].thought)
		authorPickCount[candidates[bestIndex].thought.UserID]++
	}

	return selected
}
