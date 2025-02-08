package handlers

import (
	"family-tree-app/database"
	"family-tree-app/models"
	"html/template"
	"net/http"
)

// ModerationHandler displays the moderation queue and handles approval/rejection
func ModerationHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch the admin user
	var admin models.User
	if err := database.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		http.Error(w, "Admin not found", http.StatusUnauthorized)
		return
	}

	// Check if the user is an admin
	if admin.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		// Handle approval/rejection
		editID := r.FormValue("edit_id")
		action := r.FormValue("action") // "approve" or "reject"

		var edit models.UserEdit
		if err := database.DB.First(&edit, editID).Error; err != nil {
			http.Error(w, "Edit not found", http.StatusNotFound)
			return
		}

		// Update the edit status
		if action == "approve" {
			edit.Status = "approved"
		} else if action == "reject" {
			edit.Status = "rejected"
		} else {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		// Save the updated edit
		if err := database.DB.Save(&edit).Error; err != nil {
			http.Error(w, "Failed to update edit", http.StatusInternalServerError)
			return
		}

		// Add to moderation queue
		moderationEntry := models.ModerationQueue{
			EditID:  edit.ID,
			AdminID: admin.ID,
			Status:  edit.Status,
		}
		if err := database.DB.Create(&moderationEntry).Error; err != nil {
			http.Error(w, "Failed to log moderation", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/moderation", http.StatusSeeOther)
	} else {
		// Display the moderation queue
		var pendingEdits []models.UserEdit
		if err := database.DB.Where("status = ?", "pending").Find(&pendingEdits).Error; err != nil {
			http.Error(w, "Failed to fetch pending edits", http.StatusInternalServerError)
			return
		}

		// Render the moderation queue (you can use HTML templates here)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Moderation Queue</h1>"))
		for _, edit := range pendingEdits {
			w.Write([]byte("<div>"))
			w.Write([]byte("<p>Edit ID: " + string(edit.ID) + "</p>"))
			w.Write([]byte("<p>Tree ID: " + string(edit.TreeID) + "</p>"))
			w.Write([]byte("<p>Edit Data: " + edit.EditData + "</p>"))
			w.Write([]byte("<form method='POST' action='/moderation'>"))
			w.Write([]byte("<input type='hidden' name='edit_id' value='" + string(edit.ID) + "'>"))
			w.Write([]byte("<button type='submit' name='action' value='approve'>Approve</button>"))
			w.Write([]byte("<button type='submit' name='action' value='reject'>Reject</button>"))
			w.Write([]byte("</form>"))
			w.Write([]byte("</div>"))
		}
		// Render the moderation queue using the template
		tmpl, err := template.ParseFiles("templates/moderation.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, pendingEdits)
	}
	if r.Method == "GET" {
	}
}
