package pages

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var ApprovalNeeded = false // if new posts need moderator approval

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		forPage.Tags = database.Tags

		// Check if the user is logged in
		if !forPage.User.LoggedIn {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		// Handle GET request to render the create post page
		if r.Method == http.MethodGet {
			handlers.RenderTemplates("createPost", forPage, w, r)
			return
		}

		// Handle POST request
		if r.Method != http.MethodPost {
			return
		}

		title, description, selectedTags, errStr := CheckPostCreation(r)
		if errStr != "" {
			handlers.ErrorHandler(w, http.StatusInternalServerError, errStr)
			return
		}

		// Check if an image file was uploaded
		imageFilename, errStr := ProcessImageUpload(r)
		if errStr != "" {
			forPage.Error.Error = true
			forPage.Error.Message = errStr
			forPage.Error.Field1 = title
			forPage.Error.Field2 = description
			forPage.Error.Field3 = selectedTags
			handlers.RenderTemplates("createPost", forPage, w, r)
			return
		}

		// Prepare the SQL statement for inserting post data
		stmt := "INSERT INTO posts (username, title, description, imageFilename, approved) VALUES (?, ?, ?, ?, ?)"

		var approval2 bool
		if forPage.User.TypeInt > 0 {
			approval2 = false
		} else {
			approval2 = ApprovalNeeded
		}
		// Execute the SQL statement to insert post data into the database
		result, err := db.Exec(stmt, forPage.User.Username, title, description, imageFilename, !approval2)
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to insert post data into the database.")
			return
		}

		// Get the ID of the newly created post
		postID, err := result.LastInsertId()
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to get post ID.")
			return
		}

		// Insert the selected tags into the post_tags table
		for _, tagID := range selectedTags {
			_, err = db.Exec("INSERT INTO post_tags (postID, tagID) VALUES (?, ?)", postID, tagID)
			if err != nil {
				handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to insert tag into post_tags table.")
				return
			}
		}

		// Redirect to the homepage or display a success message
		http.Redirect(w, r, "/viewPost?id="+strconv.FormatInt(postID, 10), http.StatusFound)
	}
}

func CheckPostCreation(r *http.Request) (string, string, []string, string) {
	// Parse the form data from the request
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", "", []string{}, "Failed to parse form data"
	}

	// Extract the post data
	title := r.FormValue("title")
	description := r.FormValue("description")
	selectedTags := r.Form["tags"]

	// Check for empty fields
	if title == "" || description == "" || len(selectedTags) == 0 {
		return "", "", []string{}, "Forbidden empty fields"
	}

	return title, description, selectedTags, ""
}

func ProcessImageUpload(r *http.Request) (string, string) {
	maxFileSize := 20 * 1024 * 1024 // 20MB
	file, header, err := r.FormFile("image")
	if err != nil {
		// No image file uploaded
		return "", ""
	}
	defer file.Close()

	// Validate the image size
	if header.Size > int64(maxFileSize) {
		return header.Filename, "Image size exceeds the maximum limit of 20MB"
	}

	// Validate the image format
	fileExt := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[fileExt] {
		return "", "Invalid image format. Only JPEG, PNG, and GIF formats are allowed"
	}

	// Generate a unique filename for the image
	imageFilename := generateUniqueFileName(header.Filename)

	// Save the image file to the server
	imagePath := filepath.Join("./assets/uploads", imageFilename)
	err = saveImageToFile(file, imagePath)
	if err != nil {
		return header.Filename, "Failed to save image."
	}

	return imageFilename, ""
}

func saveImageToFile(file multipart.File, imagePath string) error {
	// Create a new file at the specified path
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func generateUniqueFileName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	return filename
}
