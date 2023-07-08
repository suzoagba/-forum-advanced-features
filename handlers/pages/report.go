package pages

import (
	"database/sql"
	"forum/handlers"
	"net/http"
)

func ReportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := handlers.IsLoggedIn(r, db).User
		if r.Method != http.MethodPost || user.TypeInt != 1 {
			http.Redirect(w, r, "/", http.StatusFound)
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
		SET reported = true, report_reason = ?
		WHERE postID = ?
	`
	_, err = db.Exec(updateQuery, reason, postID)
	if err != nil {
		return err
	}
	return nil
}
