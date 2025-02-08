package handlers

import (
	"family-tree-app/database"
	"family-tree-app/models"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			// User is not authenticated, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// User is authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Fetch the user from the database
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		// User is an admin, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
