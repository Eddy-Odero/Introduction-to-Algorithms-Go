package handlers

import (
	"forum/internal/db"
	"forum/internal/middleware"
	"forum/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// PostPageData matches template structure
type PostPageData struct {
	LoggedIn     bool
	CategoryName string
	Posts        []models.Post
	SinglePost   models.Post
	Comments     []Comment
	Error        string
	Email        string
	Username     string
	UserVotes    map[int]string
	CommentVotes map[int]string
}

type Comment struct {
	ID       int
	Author   string
	Content  string
	Likes    int
	Dislikes int
}

func getUser(r *http.Request) middleware.UserContext {
	return middleware.GetUserFromContext(r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	postRepo := models.NewPostRepository(db.DB)
	posts, err := postRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	user := getUser(r)
	userVotes := make(map[int]string)
	if user.LoggedIn {
		reactionRepo := models.NewReactionRepository(db.DB)
		userVotes, _ = reactionRepo.GetUserPostVotes(int(user.ID))
	}

	data := PostPageData{
		Posts:     posts,
		LoggedIn:  user.LoggedIn,
		Username:  user.Username,
		UserVotes: userVotes,
	}

	renderTemplate(w, "index.html", data)
}

func CategoryHandler(w http.ResponseWriter, r *http.Request, categoryName string) {
	postRepo := models.NewPostRepository(db.DB)
	posts, err := postRepo.GetByCategory(categoryName)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	user := getUser(r)
	data := PostPageData{
		LoggedIn:     user.LoggedIn,
		Username:     user.Username,
		CategoryName: categoryName,
		Posts:        posts,
		UserVotes:    make(map[int]string),
	}

	renderTemplate(w, "category.html", data)
}
func PostDetailHandler(w http.ResponseWriter, r *http.Request, postID int) {
	postRepo := models.NewPostRepository(db.DB)
	post, err := postRepo.GetByID(postID)
	if err != nil {
		http.Error(w, "Failed to load post", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	commentRepo := models.NewCommentRepository(db.DB)
	modelComments, err := commentRepo.GetCommentsByPost(int64(postID))
	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}
	comments := make([]Comment, 0, len(modelComments))
	for _, c := range modelComments {
		comments = append(comments, Comment{
			ID:      int(c.ID),
			Author:  c.Username,
			Content: c.Content,
		})
	}

	user := getUser(r)
	userVotes := make(map[int]string)
	if user.LoggedIn {
		reactionRepo := models.NewReactionRepository(db.DB)
		userVotes, _ = reactionRepo.GetUserPostVotes(int(user.ID))
	}

	data := PostPageData{
		LoggedIn:   user.LoggedIn,
		Username:   user.Username,
		SinglePost: *post,
		UserVotes:  userVotes,
		Comments:   comments,
	}

	renderTemplate(w, "post-detail.html", data)
}

func NewPostPageHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	data := PostPageData{
		LoggedIn: user.LoggedIn,
		Username: user.Username,
	}

	renderTemplate(w, "new-post.html", data)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categoryInput := r.FormValue("category")

	if title == "" || content == "" {
		user := getUser(r)
		data := PostPageData{
			LoggedIn: user.LoggedIn,
			Username: user.Username,
			Error:    "Title and content are required.",
		}
		renderTemplate(w, "new-post.html", data)
		return
	}

	user := getUser(r)
	if !user.LoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postRepo := models.NewPostRepository(db.DB)
	postID, err := postRepo.CreatePost(int(user.ID), title, content)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	if categoryInput != "" {
		var categories []string
		for _, cat := range strings.Split(categoryInput, ",") {
			cleaned := strings.TrimSpace(cat)
			if cleaned != "" {
				categories = append(categories, cleaned)
			}
		}
		if len(categories) > 0 {
			catRepo := models.NewCategoryRepository(db.DB)
			catRepo.LinkToPost(int(postID), categories)
		}
	}

	http.Redirect(w, r, "/post?id="+strconv.FormatInt(postID, 10), http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, templateName string, data PostPageData) {
	uiDir := "web"
	tmpl, err := template.ParseFiles(
		filepath.Join(uiDir, "templates", "base.html"),
		filepath.Join(uiDir, "templates", templateName),
	)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base", data)
}