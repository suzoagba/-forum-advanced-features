package pages

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User

		if forPage.User.TypeInt != 2 {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		if r.Method == http.MethodGet {
			var err error
			forPage.User.Admin.Notifications, err = getAdminNotifications(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			handlers.RenderTemplates("admin", forPage, w, r)
			return
		} else {
			var err error
			page := strings.Split(r.URL.Path, "/")[2]

			if page == "post" {
				id := r.FormValue("id")
				approved := r.FormValue("approved")
				if approved == "true" {
					err = ApprovePost(db, id)
				} else {
					err = deletePost(db, id)
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				err = deleteAdminNotification(db, true, id)
			} else if page == "user" {
				id := r.FormValue("id")
				name := r.FormValue("name")
				level := r.FormValue("level")
				err = changeUserLevel(db, id, level)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				err = deleteAdminNotification(db, false, id)
				http.Redirect(w, r, "/user?name="+name, http.StatusFound)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			forPage.User = handlers.IsLoggedIn(r, db).User
			handlers.RenderTemplates("admin", forPage, w, r)
			return
		}
	}
}

func getAdminNotifications(db *sql.DB) ([]structs.AdminNotification, error) {
	selectQuery := `
		SELECT id, post, postID, user, userID
		FROM admin_notifications
	`
	rows, err := db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []structs.AdminNotification{}
	var postID sql.NullString
	var userID sql.NullString
	for rows.Next() {
		log.Println("next")
		var notification structs.AdminNotification
		err := rows.Scan(
			&notification.ID,
			&notification.Post,
			&postID,
			&notification.User,
			&userID,
		)
		if err != nil {
			return nil, err
		}
		log.Println(notification)

		if postID.Valid {
			notification.PostID, _ = strconv.Atoi(postID.String)
		} else {
			notification.PostID = 0
		}
		if userID.Valid {
			notification.UserID = userID.String
		} else {
			notification.UserID = ""
		}

		if notification.User {
			// Retrieve the usernames based on the user IDs
			notification.Username, err = handlers.GetUsernameByID(db, notification.UserID)
			if err != nil {
				return nil, err
			}
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Println(notifications)

	/*	err = MarkNotificationsAsRead(db, userID)

		if err = rows.Err(); err != nil {
			return nil, err
		}*/

	return notifications, nil
}

func deleteAdminNotification(db *sql.DB, post bool, id string) error {
	var what string
	if post {
		what = "postID"
	} else {
		what = "userID"
	}

	deleteQuery := `
		DELETE FROM admin_notifications
		WHERE post = ? AND ` + what + ` = ?
	`
	_, err := db.Exec(deleteQuery, post, id)
	if err != nil {
		return err
	}
	return nil
}

func changeUserLevel(db *sql.DB, id string, level string) error {
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Update the user's level
	updateUserLevelQuery := `
		UPDATE users
		SET level = ?, requested_for_promotion = false
		WHERE uuid = ?
	`
	_, err = tx.Exec(updateUserLevelQuery, level, id)
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

	return nil
}
