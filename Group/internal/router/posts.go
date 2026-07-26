package router

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"forum/internal/models"
)

func RegisterPostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /post", postDetailHandler)
	mux.HandleFunc("GET /category/{name}", categoryHandler)

	mux.HandleFunc("GET /new-post", middleware.RequireAuth(handlers.NewPostPageHandler))
	mux.HandleFunc("POST /new-post", middleware.RequireAuth(handlers.CreatePostHandler))

	// Profile with filters (Phase 8)
	mux.HandleFunc("GET /profile", middleware.RequireAuth(profileHandler))
}

func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.CustomNotFoundHandler(w, r)
		return
	}
	handlers.PostDetailHandler(w, r, postID)
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.PathValue("name")
	if categoryName == "" {
		middleware.CustomNotFoundHandler(w, r)
		return
	}
	handlers.CategoryHandler(w, r, categoryName)
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)

	filter := strings.TrimSpace(r.URL.Query().Get("filter"))
	postRepo := models.NewPostRepository(db.DB)

	var posts []models.Post
	var err error
	switch filter {
	case "liked":
		posts, err = postRepo.GetLikedByUserID(int(user.ID))
	case "mine":
		fallthrough
	default:
		filter = "mine"
		posts, err = postRepo.GetByUserID(int(user.ID))
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get actual comment count
	commentRepo := models.NewCommentRepository(db.DB)
	commentCount, _ := commentRepo.CountCommentsByUserID(int64(user.ID))
	reactionRepo := models.NewReactionRepository(db.DB)
	likesGiven, _ := reactionRepo.CountLikesGivenByUserID(int(user.ID))

	// Convert models.Post to router.Post with ALL fields
	var pagePosts []Post
	for _, p := range posts {
		pagePosts = append(pagePosts, Post{
			ID:           p.ID,
			Title:        p.Title,
			Author:       p.Author,
			Excerpt:      p.Excerpt,
			Content:      p.Content,
			ImageURL:     p.ImageURL,
			Categories:   p.Categories,
			Likes:        p.Likes,
			Dislikes:     p.Dislikes,
			CommentCount: p.CommentCount,
		})
	}
data := PageData{
		LoggedIn:      true,
		Username:      user.Username,
		CurrentFilter: filter,
		Posts:         pagePosts,
		PostCount:     len(pagePosts),
		CommentCount:  int(commentCount),
		LikesGiven:    likesGiven,
	}

	RenderTemplate(w, "profile.html", data)
}