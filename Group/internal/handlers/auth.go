package handlers

import (
	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/middleware"
	"forum/internal/models"
	"net/http"
	"strings"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	data := PostPageData{
		LoggedIn: user.LoggedIn,
		Username: user.Username,
	}
	renderTemplate(w, "register.html", data)
}

// RegisterHandler processes registration form submissions
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	// Validation
	if email == "" || username == "" || password == "" {
		data := PostPageData{
			LoggedIn: false,
			Error:    "All fields are required.",
		}
		renderTemplate(w, "register.html", data)
		return
	}

	userRepo := models.NewUserRepository(db.DB)

	// Check if email already exists
	exists, _ := userRepo.EmailExists(email)
	if exists {
		data := PostPageData{
			LoggedIn: false,
			Email:    email,
			Username: username,
			Error:    "Email already registered.",
		}
		renderTemplate(w, "register.html", data)
		return
	}

	// Check if username already exists
	exists, _ = userRepo.UsernameExists(username)
	if exists {
		data := PostPageData{
			LoggedIn: false,
			Email:    email,
			Username: username,
			Error:    "Username already taken.",
		}
		renderTemplate(w, "register.html", data)
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		data := PostPageData{
			LoggedIn: false,
			Error:    "Internal error. Please try again.",
		}
		renderTemplate(w, "register.html", data)
		return
	}

	// Create user
	_, err = userRepo.CreateUser(username, email, passwordHash)
	if err != nil {
		data := PostPageData{
			LoggedIn: false,
			Error:    "Failed to create account. Please try again.",
		}
		renderTemplate(w, "register.html", data)
		return
	}

	// Redirect to login on success
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LoginPageHandler renders the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	data := PostPageData{
		LoggedIn: user.LoggedIn,
		Username: user.Username,
	}
	renderTemplate(w, "login.html", data)
}

// LoginHandler processes login form submissions
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		data := PostPageData{
			LoggedIn: false,
			Error:    "Email and password are required.",
		}
		renderTemplate(w, "login.html", data)
		return
	}

	userRepo := models.NewUserRepository(db.DB)
	user, err := userRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		data := PostPageData{
			LoggedIn: false,
			Email:    email,
			Error:    "Invalid email or password.",
		}
		renderTemplate(w, "login.html", data)
		return
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		data := PostPageData{
			LoggedIn: false,
			Email:    email,
			Error:    "Invalid email or password.",
		}
		renderTemplate(w, "login.html", data)
		return
	}

	// Create session
	sessionRepo := auth.NewSessionRepository(db.DB)
	token, err := sessionRepo.CreateSession(int(user.ID))
	if err != nil {
		data := PostPageData{
			LoggedIn: false,
			Error:    "Failed to create session. Please try again.",
		}
		renderTemplate(w, "login.html", data)
		return
	}

	// Set cookie
	auth.SetSessionCookie(w, token)

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler clears the session and redirects to homepage
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := auth.GetSessionTokenFromCookie(r)
	if token != "" {
		sessionRepo := auth.NewSessionRepository(db.DB)
		sessionRepo.DeleteSession(token)
	}

	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
