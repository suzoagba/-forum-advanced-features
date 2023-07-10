package handlers

import (
	"database/sql"
	"forum/structs"
)

func CreateNotification(db *sql.DB, action, who string, isPost, isComment bool, postID, commentID string) error {
	var userID string
	var err error
	if isPost {
		userID, err = GetUserIDFromPostID(db, postID)
		if err != nil {
			return err
		}
	} else {
		userID, err = GetUserIDFromCommentID(db, commentID)
		if err != nil {
			return err
		}
	}

	if userID == who {
		return nil
	}

	insertQuery := `
		INSERT INTO notifications (userID, whoID, actionDone, isPost, isComment, postID, commentID)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(insertQuery, userID, who, action, isPost, isComment, postID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func MarkNotificationsAsRead(db *sql.DB, userID string) error {
	updateQuery := `
		UPDATE notifications
		SET isRead = true
		WHERE userID = ?
	`
	_, err := db.Exec(updateQuery, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNotifications(db *sql.DB, userID string) ([]structs.Notification, error) {
	selectQuery := `
		SELECT notificationID, userID, whoID, actionDone, isPost, isComment, postID, commentID, isRead, createdDate
		FROM notifications
		WHERE userID = ?
		ORDER BY createdDate DESC
	`
	rows, err := db.Query(selectQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []structs.Notification{}
	for rows.Next() {
		var notification structs.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.User,
			&notification.Who,
			&notification.ActionDone,
			&notification.IsPost,
			&notification.IsComment,
			&notification.PostID,
			&notification.CommentID,
			&notification.IsRead,
			&notification.CreatedDate,
		)
		if err != nil {
			return nil, err
		}

		// Retrieve the usernames based on the user IDs
		notification.User, err = GetUsernameByID(db, notification.User)
		if err != nil {
			return nil, err
		}

		notification.Who, err = GetUsernameByID(db, notification.Who)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	/*	err = MarkNotificationsAsRead(db, userID)

		if err = rows.Err(); err != nil {
			return nil, err
		}*/

	return notifications, nil
}

func GetUnreadNotificationCount(db *sql.DB, userID string, level int) (int, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE userID = ? AND isRead = false
	`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if level == 1 {
		mode, err := GetModeratorNotificationCount(db)

		if err != nil {
			return 0, err
		}

		count += mode
	}

	return count, nil
}

func GetUsernameByID(db *sql.DB, userID string) (string, error) {
	query := `
		SELECT username FROM users WHERE uuid = ?
	`
	row := db.QueryRow(query, userID)

	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func GetUserIDFromPostID(db *sql.DB, postID string) (string, error) {
	query := `
		SELECT username FROM posts WHERE postID = ?
	`
	row := db.QueryRow(query, postID)

	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}

	query = `
		SELECT uuid FROM users WHERE username = ?
	`
	row = db.QueryRow(query, username)

	var userID string
	err = row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetUserIDFromCommentID(db *sql.DB, commentID string) (string, error) {
	query := `
		SELECT username FROM comments WHERE commentID = ?
	`
	row := db.QueryRow(query, commentID)

	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}

	query = `
		SELECT uuid FROM users WHERE username = ?
	`
	row = db.QueryRow(query, username)

	var userID string
	err = row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetAdminNotificationCount(db *sql.DB) (int, error) {
	query := `
		SELECT COUNT(*) FROM admin_notifications
	`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetModeratorNotificationCount(db *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM posts WHERE approved = false AND reported = false"

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetModeratorPosts(db *sql.DB) ([]structs.AdminNotification, error) {
	query := `
		SELECT p.postID, u.username
		FROM posts p
		JOIN users u ON p.username = u.username
		WHERE p.approved = false AND p.reported = false
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.AdminNotification

	for rows.Next() {
		var post structs.AdminNotification

		err := rows.Scan(&post.ID, &post.Username)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
