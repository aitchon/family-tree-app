package handlers

import (
	"family-tree-app/database"
	"family-tree-app/models"
	"net/http"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		user.Username = username
		user.Role = "user" // Set default role

		if err := user.HashPassword(password); err != nil {
			http.Error(w, "Unable to hash password", http.StatusInternalServerError)
			return
		}

		if err := database.DB.Create(&user).Error; err != nil {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// Render registration form (you can use HTML templates here)
		http.ServeFile(w, r, "templates/register.html")
	}
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := user.CheckPassword(password); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Create a session
		session, _ := Store.Get(r, "session-name")
		session.Values["username"] = user.Username
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		// Render login form (you can use HTML templates here)
		http.ServeFile(w, r, "templates/login.html")
	}
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	session.Values["username"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
