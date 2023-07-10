package pages

import (
	"database/sql"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"log"
	"net/http"
)

func DeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := handlers.IsLoggedIn(r, db).User
		type forPage struct {
			structs.User
			Post    structs.Post
			Comment structs.Comment
			Tags    []structs.Tag
			Error   structs.ErrorMessage
		}

		var (
			id       string
			editType string
			ans      string
		)

		if r.Method == http.MethodGet {
			id = r.URL.Query().Get("id")
			editType = r.URL.Query().Get("type")
		} else if r.Method == http.MethodPost {
			id = r.FormValue("id")
			editType = r.FormValue("type")
			ans = r.FormValue("yesno")
		}

		var post structs.Post
		var comment structs.Comment
		if editType == "post" {
			post, _ = GetPost(w, db, id)
			if user.TypeInt == 0 {
				if post.Username != user.Username {
					handlers.ErrorHandler(w, http.StatusUnauthorized, "You are not allowed to delete other users posts")
					return
				}
			}
		} else if editType == "comment" {
			comment, _ = GetComment(w, db, id)
			if user.TypeInt != 2 {
				if comment.Username != user.Username {
					handlers.ErrorHandler(w, http.StatusUnauthorized, "You are not allowed to delete other users comments")
					return
				}
			}
		}
		if r.Method == http.MethodGet {

			data := forPage{
				User:    user,
				Post:    post,
				Comment: comment,
				Tags:    database.Tags,
			}
			handlers.RenderTemplates("delete", data, w, r)

		} else if r.Method == http.MethodPost {

			if ans == "true" {
				var err error
				if editType == "post" {
					err = deletePost(db, id)
				} else if editType == "comment" {
					err = deleteComment(db, id)
				}
				if err != nil {
					http.Error(w, "Failed to delete the "+editType, http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, "/activity", http.StatusFound)
		}
	}
}

func deletePost(db *sql.DB, postID string) error {
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete post reactions
	deletePostReactionsQuery := `
		DELETE FROM post_reactions
		WHERE post_id = ?
	`
	_, err = tx.Exec(deletePostReactionsQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete comment reactions
	deleteCommentReactionsQuery := `
		DELETE FROM comment_reactions
		WHERE comment_id IN (
			SELECT commentID FROM comments WHERE postID = ?
		)
	`
	_, err = tx.Exec(deleteCommentReactionsQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete comments
	deleteCommentsQuery := `
		DELETE FROM comments
		WHERE postID = ?
	`
	_, err = tx.Exec(deleteCommentsQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete post tags
	deletePostTagsQuery := `
		DELETE FROM post_tags
		WHERE postID = ?
	`
	_, err = tx.Exec(deletePostTagsQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the post
	deletePostQuery := `
		DELETE FROM posts
		WHERE postID = ?
	`
	_, err = tx.Exec(deletePostQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete post notifications
	deletePostNotificationsQuery := `
		DELETE FROM notifications
		WHERE postID = ?
	`
	_, err = tx.Exec(deletePostNotificationsQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	log.Printf("Post with ID %s deleted successfully", postID)
	return nil
}

func deleteComment(db *sql.DB, commentID string) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete the comment reactions associated with the comment
	_, err = tx.Exec("DELETE FROM comment_reactions WHERE comment_id = ?", commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the comment itself
	_, err = tx.Exec("DELETE FROM comments WHERE commentID = ?", commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete comment notifications
	deleteCommentNotificationsQuery := `
		DELETE FROM notifications
		WHERE commentID = ?
	`
	_, err = tx.Exec(deleteCommentNotificationsQuery, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	log.Printf("Comment with ID %s deleted successfully", commentID)
	return nil
}
