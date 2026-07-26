package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/db"
	"forum/internal/middleware"
	"forum/internal/models"
)

func CreateReactionHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if !user.LoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	targetIDStr := r.FormValue("target_id")
	targetType := r.FormValue("type")
	value := r.FormValue("val")

	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil || (value != "like" && value != "dislike") {
		http.Error(w, "Invalid reaction", http.StatusBadRequest)
		return
	}

	reactionRepo := models.NewReactionRepository(db.DB)
	switch targetType {
	case "post":
		err = reactionRepo.SetPostReaction(int(user.ID), targetID, value)
	case "comment":
		err = reactionRepo.SetCommentReaction(int(user.ID), targetID, value)
	default:
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}