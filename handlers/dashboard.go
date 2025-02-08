package handlers

import (
	"family-tree-app/database"
	"family-tree-app/models"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current user from the session
	session, _ := Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Fetch the user's uploaded GEDCOM files
	var gedcomFiles []models.GEDCOMFile
	if err := database.DB.Where("user_id = ?", user.ID).Find(&gedcomFiles).Error; err != nil {
		http.Error(w, "Failed to fetch GEDCOM files", http.StatusInternalServerError)
		return
	}

	// Render the dashboard (you can use HTML templates here)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Dashboard</h1>"))
	w.Write([]byte("<h2>Your Family Trees</h2>"))
	for _, file := range gedcomFiles {
		w.Write([]byte("<div>"))
		w.Write([]byte("<p>Filename: " + file.Filename + "</p>"))
		w.Write([]byte("<a href='/tree/" + string(file.ID) + "'>View Tree</a>"))
		w.Write([]byte("</div>"))
	}
	w.Write([]byte("<a href='/upload'>Upload New GEDCOM File</a>"))
}
