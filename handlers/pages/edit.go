package pages

import (
	"database/sql"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func EditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := handlers.IsLoggedIn(r, db).User
		type forPage struct {
			structs.User
			Post    structs.Post
			Comment structs.Comment
			Tags    []structs.Tag
		}

		id := r.URL.Query().Get("id")
		editType := r.URL.Query().Get("type")

		var post structs.Post
		var comment structs.Comment
		if r.Method == http.MethodGet {
			if editType == "post" {
				post, _ = GetPost(w, db, id)
			} else if editType == "comment" {
				comment, _ = GetComment(w, db, id)
			}

			// Create a data struct to pass to the template
			data := forPage{
				User:    user,
				Post:    post,
				Comment: comment,
				Tags:    database.Tags,
			}
			handlers.RenderTemplates("edit", data, w, r)
		} else if r.Method == http.MethodPost {
			//handlers.RenderTemplates("edit", data, w, r) // TODO
		}
	}
}

// GetPost from database
func GetPost(w http.ResponseWriter, db *sql.DB, postID string) (structs.Post, error) {
	postQuery := `
			SELECT postID, title, description, imageFileName, creationDate, username, likes, dislikes
			FROM posts
			WHERE postID = ?
		`

	postRow := db.QueryRow(postQuery, postID)

	var post structs.Post
	var imageFileName sql.NullString
	err := postRow.Scan(&post.ID, &post.Title, &post.Description, &imageFileName, &post.CreationDate, &post.Username, &post.Likes, &post.Dislikes)
	if err != nil {
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return structs.Post{}, err
	}

	if imageFileName.Valid {
		post.ImageFileName = imageFileName.String
	} else {
		post.ImageFileName = "" // Set a default value for imageFileName when it is NULL
	}

	post.Tags, _ = database.GetPostTags(db, postID)

	return post, nil
}

func GetComment(w http.ResponseWriter, db *sql.DB, commentID string) (structs.Comment, error) {
	commentQuery := `
			SELECT commentID, postID, content, username, creationDate, likes, dislikes
			FROM comments
			WHERE commentID = ?
		`

	commentRow := db.QueryRow(commentQuery, commentID)

	var comment structs.Comment
	err := commentRow.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.Username, &comment.CreationDate, &comment.Likes, &comment.Dislikes)
	if err != nil {
		http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
		return structs.Comment{}, err
	}

	return comment, nil
}
