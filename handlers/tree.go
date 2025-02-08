package handlers

import (
	"family-tree-app/database"
	"family-tree-app/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func TreeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the tree ID from the URL
	vars := mux.Vars(r)
	treeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid tree ID", http.StatusBadRequest)
		return
	}

	// Fetch the family tree from the database
	var familyTree models.FamilyTree
	if err := database.DB.First(&familyTree, treeID).Error; err != nil {
		http.Error(w, "Family tree not found", http.StatusNotFound)
		return
	}

	// Return the family tree data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(familyTree.Data))

	// Serve the HTML page
	http.ServeFile(w, r, "templates/tree.html")
}
