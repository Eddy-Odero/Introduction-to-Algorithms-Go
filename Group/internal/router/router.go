package router

import (
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"forum/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

type Post struct {
	ID           int
	Title        string
	Author       string
	Excerpt      string
	Content      string
	ImageURL     string
	Categories   []string
	Likes        int
	Dislikes     int
	CommentCount int
}

type Comment struct {
	ID       int
	Author   string
	Content  string
	Likes    int
	Dislikes int
}

type PageData struct {
	LoggedIn      bool
	CategoryName  string
	CurrentFilter string
	Posts         []Post
	SinglePost    Post
	Comments      []Comment
	PostCount     int
	CommentCount  int
	LikesGiven    int
	Error         string
	Email         string
	Username      string
	UserVotes     map[int]string
	CommentVotes  map[int]string
}

var uiDir string

func SetupRouter(routerUIDir string) *http.ServeMux {
	mux := http.NewServeMux()
	uiDir = routerUIDir

	// Static files
	fileServer := http.FileServer(http.Dir(routerUIDir + "/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Phase 5: All post routes
	RegisterPostRoutes(mux)

	// Phase 4: Auth routes (real session-based)
	mux.HandleFunc("GET /login", handlers.LoginPageHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("GET /register", handlers.RegisterPageHandler)
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("GET /logout", handlers.LogoutHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

	// Phase 6: Comment routes
	commentHandler := handlers.NewCommentHandler(models.NewCommentRepository(db.DB))
	mux.HandleFunc("POST /comment", middleware.RequireAuth(commentHandler.CreateComment))
    mux.HandleFunc("POST /react", middleware.RequireAuth(handlers.CreateReactionHandler))

	// Catch-all for 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			middleware.CustomNotFoundHandler(w, r)
			return
		}
	})

	return mux
}

func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
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