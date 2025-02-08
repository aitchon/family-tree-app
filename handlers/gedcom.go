package handlers

import (
	"encoding/json"
	"family-tree-app/database"
	"family-tree-app/models"
	"io"
	"net/http"

	"github.com/elliotchance/gedcom"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, fileHeader, err := r.FormFile("gedcom")
		if err != nil {
			http.Error(w, "Unable to upload file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get the filename
		filename := fileHeader.Filename

		// Read the file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file content", http.StatusInternalServerError)
			return
		}

		// Parse GEDCOM file from the uploaded file
		document, err := gedcom.NewDocumentFromGEDCOMFile(filename)
		if err != nil {
			http.Error(w, "Invalid GEDCOM file", http.StatusBadRequest)
			return
		}

		// Extract family tree data from the GEDCOM document
		familyTreeData := extractFamilyTreeData(document)

		// Convert the family tree data to JSON
		treeDataJSON, err := json.Marshal(familyTreeData)
		if err != nil {
			http.Error(w, "Failed to convert family tree data to JSON", http.StatusInternalServerError)
			return
		}

		// Get the current user from the session
		session, _ := Store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Save the GEDCOM file
		gedcomFile := models.GEDCOMFile{
			UserID:   user.ID,
			Filename: filename, // Use the filename from the file header
			Content:  fileContent,
		}
		if err := database.DB.Create(&gedcomFile).Error; err != nil {
			http.Error(w, "Failed to save GEDCOM file", http.StatusInternalServerError)
			return
		}

		// Save the family tree
		familyTree := models.FamilyTree{
			GEDCOMFileID: gedcomFile.ID,
			Data:         string(treeDataJSON), // Store the JSON data
		}
		if err := database.DB.Create(&familyTree).Error; err != nil {
			http.Error(w, "Failed to save family tree", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/tree/"+string(familyTree.ID), http.StatusSeeOther)
	} else {
		// Render upload form
		http.ServeFile(w, r, "templates/upload.html")
	}
}

// extractFamilyTreeData extracts family tree data from a GEDCOM document
func extractFamilyTreeData(document *gedcom.Document) map[string]interface{} {
	// Create a map to store individuals by their ID
	individualsMap := make(map[string]map[string]interface{})
	for _, individual := range document.Individuals() {
		birthDate, _ := individual.Birth()
		deathDate, _ := individual.Death()

		individualsMap[individual.Pointer()] = map[string]interface{}{
			"id":         individual.Pointer(),
			"name":       individual.Name().String(),
			"birth_date": birthDate.String(),
			"death_date": deathDate.String(),
			"children":   []map[string]interface{}{},
		}
	}

	// Build the hierarchical structure
	for _, family := range document.Families() {
		husbandID := family.Husband().Pointer()
		wifeID := family.Wife().Pointer()

		for _, child := range family.Children() {
			childID := child.Individual().Pointer()
			individualsMap[husbandID]["children"] = append(individualsMap[husbandID]["children"].([]map[string]interface{}), individualsMap[childID])
			individualsMap[wifeID]["children"] = append(individualsMap[wifeID]["children"].([]map[string]interface{}), individualsMap[childID])
		}
	}

	// Find the root individual (e.g., the first individual in the map)
	for _, individual := range individualsMap {
		return individual
	}

	return nil
}

// getChildrenIDs extracts the IDs of children in a family
func getChildrenIDs(children []*gedcom.ChildNode) []string {
	ids := make([]string, 0)
	for _, child := range children {
		ids = append(ids, child.Individual().Pointer()) // Use Pointer() for the unique identifier
	}
	return ids
}
