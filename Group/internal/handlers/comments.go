package handlers

import (
	"forum/internal/middleware"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

// UserIDExtractor defines the function signature for getting the current user.
// This allows dependency injection for testing.
type UserIDExtractor func(r *http.Request) (int64, bool)

// CommentHandler processes comment-related HTTP requests.
type CommentHandler struct {
	CommentRepo *models.CommentRepository
	GetUserID   UserIDExtractor
}

// NewCommentHandler creates a handler with the given repository.
func NewCommentHandler(repo *models.CommentRepository) *CommentHandler {
	return &CommentHandler{
		CommentRepo: repo,
		GetUserID:   defaultGetUserID,
	}
}

// CreateComment handles POST /comment.
// Authenticated users submit comments; guests are redirected to login.
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := h.GetUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = h.CommentRepo.CreateComment(postID, userID, content)
	if err != nil {
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.FormatInt(postID, 10), http.StatusSeeOther)
}

// defaultGetUserID obtains the user added by Phase 4 authentication middleware.
func defaultGetUserID(r *http.Request) (int64, bool) {
	user := middleware.GetUserFromContext(r)
	if !user.LoggedIn || user.ID <= 0 {
		return 0, false
	}
	return int64(user.ID), true
}
