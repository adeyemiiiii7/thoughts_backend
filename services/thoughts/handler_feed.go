package thoughts

import (
	"math/rand"
	"net/http"
	"sort"
	"strconv"
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

	limit := parseFeedLimit(r.URL.Query().Get("limit"))

	var followedIDs []uint
	if err := h.db.Model(&models.Follow{}).
		Where("follower_id = ?", user.ID).
		Pluck("following_id", &followedIDs).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load follow graph",
		})
		return
	}

	// Keep own thoughts in the following feed too, so the home feed never feels empty.
	authorIDs := append([]uint{user.ID}, followedIDs...)

	var thoughts []models.Thought
	if err := h.db.
		Where("user_id IN ?", authorIDs).
		Preload("User").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
		Limit(limit * 3).
		Find(&thoughts).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load following feed",
		})
		return
	}

	// Following feed is mostly chronological, but still lightly diversify authors
	// so one person does not dominate the first screen.
	candidates := make([]feedCandidate, 0, len(thoughts))
	for _, thought := range thoughts {
		score := recencyScore(thought.CreatedAt) + float64(len(thought.Comments)*2+len(thought.Reactions))
		if thought.UserID == user.ID {
			score += 40
		}

		candidates = append(candidates, feedCandidate{
			thought: thought,
			score:   score,
		})
	}

	shared.RespondJSON(w, http.StatusOK, diversifyFeed(candidates, limit))
}

func (h *Handler) FYPFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	limit := parseFeedLimit(r.URL.Query().Get("limit"))
	candidateLimit := limit * 5
	if candidateLimit < 60 {
		candidateLimit = 60
	}
	if candidateLimit > 180 {
		candidateLimit = 180
	}

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

	var thoughts []models.Thought
	if err := h.db.
		Preload("User").
		Preload("User.Interests").
		Preload("Comments").
		Preload("Reactions").
		Order("created_at DESC").
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

	// FYP should feel more exploratory than the following feed.
	// So:
	// - exclude your own posts
	// - give a small penalty to accounts you already follow
	// - boost shared interests
	// - add a little randomness so the page does not feel static
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	candidates := make([]feedCandidate, 0, len(thoughts))
	for _, thought := range thoughts {
		if thought.UserID == viewer.ID {
			continue
		}

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

		candidates = append(candidates, feedCandidate{
			thought: thought,
			score:   score,
		})
	}

	shared.RespondJSON(w, http.StatusOK, diversifyFeed(candidates, limit))
}

func parseFeedLimit(value string) int {
	if value == "" {
		return 20
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}

	return limit
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
