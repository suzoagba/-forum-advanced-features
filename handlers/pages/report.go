package pages

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func ReportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := handlers.IsLoggedIn(r, db).User
		type forPage struct {
			User  structs.User
			Post  structs.Post
			Error structs.ErrorMessage
		}
		if r.Method != http.MethodPost {
			if user.TypeInt != 1 {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			} else {
				id := r.URL.Query().Get("id")
				rType := r.URL.Query().Get("type")
				var post structs.Post
				if rType == "post" {
					post, _ = GetPost(w, db, id)
					if user.TypeInt != 1 {
						if user.TypeInt == 0 || (user.TypeInt == 1 && post.Approved) {
							handlers.ErrorHandler(w, http.StatusUnauthorized, "You are not allowed to delete other users posts")
							return
						}
					}
				}
				data := forPage{
					User: user,
					Post: post,
				}
				handlers.RenderTemplates("report", data, w, r)
				return
			}
		} else {
			id := r.FormValue("id")
			approved := r.FormValue("approved")
			if approved == "true" {
				err := ApprovePost(db, id)
				if err != nil {
					http.Error(w, "Failed to approve post", http.StatusInternalServerError)
				}
				http.Redirect(w, r, "/viewPost?id="+id, http.StatusFound)
				return
			} else {
				reason := r.FormValue("reason")
				err := storeReport(db, id, reason)
				if err != nil {
					http.Error(w, "Failed to save report", http.StatusInternalServerError)
				}
				http.Redirect(w, r, "/viewPost?id="+id, http.StatusFound)
				return
			}
		}
	}
}

func ApprovePost(db *sql.DB, id string) error {
	updateQuery := `
		UPDATE posts
		SET approved = true, reported = false
		WHERE postID = ?
	`
	_, err := db.Exec(updateQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func storeReport(db *sql.DB, postID string, reason string) error {
	insertQuery := `
		INSERT INTO admin_notifications (post, postID)
		VALUES (?, ?)
	`
	_, err := db.Exec(insertQuery, true, postID)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE posts
		SET approved = false, reported = true, report_reason = ?
		WHERE postID = ?
	`
	_, err = db.Exec(updateQuery, reason, postID)
	if err != nil {
		return err
	}
	return nil
}
