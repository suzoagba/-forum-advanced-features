package pages

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func ViewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the postID from the URL parameters
		postID := r.URL.Query().Get("id")
		react := r.URL.Query().Get("react")
		like := r.URL.Query().Get("like")
		user := handlers.IsLoggedIn(r, db).User

		// Check if it's a reaction to the post
		if react != "" {
			var ifPost bool
			var id string
			if react == "0" {
				ifPost = true
				id = postID
			} else {
				ifPost = false
				id = react
			}
			// Check if the user has already liked/disliked it
			reactionExists, wasLike, err := checkReactionExists(db, id, r, ifPost)
			if err != nil {
				http.Error(w, "Failed to check reaction existence", http.StatusInternalServerError)
				return
			}

			var errStr string
			if !reactionExists {
				// User has not reacted to it, allow like/dislike action
				errStr = addLike(db, ifPost, postID, react, user.ID, like)
			} else {
				// Update the like and dislike count for it
				errStr = updateLike(db, ifPost, wasLike, postID, react, user.ID, like)
			}
			if errStr != "" {
				http.Error(w, errStr, http.StatusInternalServerError)
				return
			}
		}

		post, _ := database.GetPost(db, postID)
		comments, _ := database.GetPostComments(db, postID)

		// Create a data struct to pass to the template
		data := struct {
			structs.User
			Post     structs.Post
			Comments []structs.Comment
		}{
			User:     handlers.IsLoggedIn(r, db).User,
			Post:     post,
			Comments: comments,
		}

		// Render the template with the data
		handlers.RenderTemplates("viewPost", data, w, r)
	}
}

func addLike(db *sql.DB, post bool, postID, commentID, userID, like string) string {
	var (
		typeOf string
		id     string
		action string
	)
	if post {
		typeOf = "post"
		id = postID
	} else {
		typeOf = "comment"
		id = commentID
	}
	if like == "true" {
		action = "like"
	} else if like == "false" {
		action = "dislike"
	}
	updateTopicQuery := "UPDATE " + typeOf + "s SET " + action + "s = " + action + "s + 1 WHERE " + typeOf + "ID = ?"
	_, err := db.Exec(updateTopicQuery, id)
	if err != nil {
		return "Failed to update " + typeOf
	}

	insertReactionQuery := `
		INSERT INTO ` + typeOf + `_reactions (` + typeOf + `_id, user_id, reaction_type)
		VALUES (?, ?, ?)
	`
	_, err = db.Exec(insertReactionQuery, id, userID, like == "true")
	if err != nil {
		return "Failed to store reaction"
	}

	err = handlers.CreateNotification(db, action, userID, true, false, postID, commentID)
	if err != nil {
		return "Failed to store notification"
	}
	return ""
}

func updateLike(db *sql.DB, post, wasLike bool, postID, commentID, userID, like string) string {
	// Update the like and dislike count for it
	var (
		typeOf           string
		likeIncrement    int
		dislikeIncrement int
		err              error
		action           string
		id               string
	)

	if post {
		typeOf = "post"
		id = postID
	} else {
		typeOf = "comment"
		id = commentID
	}
	updateTopicQuery := "UPDATE " + typeOf + "s SET likes = likes + ?, dislikes = dislikes + ? WHERE " + typeOf + "ID = ?"

	if (like == "true" && !wasLike) || (like == "false" && wasLike) {
		if like == "true" && !wasLike {
			likeIncrement = 1
			dislikeIncrement = -1
			action = "like"
		} else if like == "false" && wasLike {
			likeIncrement = -1
			dislikeIncrement = 1
			action = "dislike"
		}

		// Update the reaction_type to false for the specified topic and user
		updateQuery := "UPDATE " + typeOf + "_reactions SET reaction_type = ? WHERE " + typeOf + "_id = ? AND user_id = ?"
		_, err = db.Exec(updateQuery, like == "true", id, userID)
		if err != nil {
			return "Failed to update reaction"
		}
	} else {
		if like == "true" {
			likeIncrement = -1
			dislikeIncrement = 0
			action = "unlike"
		} else {
			likeIncrement = 0
			dislikeIncrement = -1
			action = "undislike"
		}

		// Delete the row from post_reactions table for the specified topic and user
		deleteQuery := "DELETE FROM " + typeOf + "_reactions WHERE " + typeOf + "_id = ? AND user_id = ?"
		_, err := db.Exec(deleteQuery, id, userID)
		if err != nil {
			return "Failed to delete reaction"
		}
	}

	_, err = db.Exec(updateTopicQuery, likeIncrement, dislikeIncrement, id)
	if err != nil {
		return "Failed to update " + typeOf
	}

	err = handlers.CreateNotification(db, action, userID, post, !post, postID, commentID)
	if err != nil {
		return "Failed to store notification"
	}

	return ""
}

func checkReactionExists(db *sql.DB, reactionID string, r *http.Request, isPost bool) (bool, bool, error) {
	user := handlers.IsLoggedIn(r, db).User

	// Determine the table name and column names based on whether it's a post or comment reaction
	tableName := "post_reactions"
	idColumnName := "post_id"
	userIDColumnName := "user_id"
	if !isPost {
		tableName = "comment_reactions"
		idColumnName = "comment_id"
	}

	// Query the database to check if the user has already reacted to the post or comment
	query := fmt.Sprintf(
		"SELECT EXISTS (SELECT 1 FROM %s WHERE %s = ? AND %s = ?), reaction_type FROM %s WHERE %s = ? AND %s = ?",
		tableName, idColumnName, userIDColumnName, tableName, idColumnName, userIDColumnName,
	)

	row := db.QueryRow(query, reactionID, user.ID, reactionID, user.ID)

	var exists bool
	var reactionType bool
	err := row.Scan(&exists, &reactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, indicating the reaction doesn't exist
			return false, false, nil
		}
		return false, false, err
	}

	return exists, reactionType, nil
}
