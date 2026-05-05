package thoughts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"thoughts_backend_api/models"
	"thoughts_backend_api/shared"

	"github.com/go-chi/chi/v5"
)

func (h *Handler)Update(w http.ResponseWriter, r *http.Request){
	user , ok := shared.GetUserFromContext(r.Context())
	if !ok {
		shared.RespondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized"})
		return
	}
	//get thought id from URL
	thoughtIDParam := chi.URLParam(r, "id")
	thoughtID, err := strconv.ParseUint(thoughtIDParam, 10, 64)
	if err != nil {
			shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid thought id"})
		return
	}     		

	//parse it
	var req updateThoughtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body"})
		return
	}
		//trim title and content
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	//validate both are not empty
	if req.Title == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required"})
		return
	}
	if req.Content == "" {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "content is required"})
		return
	}
	//load the thought from DB
	var thought models.Thought
	if err := h.db.First(&thought, uint(thoughtID)).Error; err != nil {
		shared.RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "thought not found"})
		return
	}
	//check if the logged in user is the owner of the thought
	if thought.UserID != user.ID {
		shared.RespondJSON(w, http.StatusForbidden, map[string]string{
			"error": "you can only update your own thoughts"})
		return
	}
	//update Title and Content
	thought.Title = req.Title
	thought.Content = req.Content

// save with GORM
	if err := h.db.Save(&thought).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to update thought"})
		return
	}
// preload User if you want a nice response
	if err := h.db.Preload("User").First(&thought, thought.ID).Error; err != nil {
		shared.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load updated thought"})
		return
	}
// return updated thought JSON
	shared.RespondJSON(w, http.StatusOK, thought)
}


