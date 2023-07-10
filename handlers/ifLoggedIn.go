package handlers

import (
	"database/sql"
	"forum/structs"
	"net/http"
)

type userInfo struct {
	User structs.User
}

func IsLoggedIn(r *http.Request, db *sql.DB) userInfo {
	info := userInfo{}

	cookie, err := r.Cookie("forum-session")
	if err != nil {
		info.User.LoggedIn = false
		// Session ID cookie not found
		return info
	}
	sessionID := cookie.Value
	if sessionID == "" {
		info.User.LoggedIn = false
		return info
	}

	// Check if the session ID exists in the database
	row := db.QueryRow(`
	SELECT u.username, u.requested_for_promotion
	FROM authenticated_users AS a
	INNER JOIN users AS u ON a.username = u.username
	WHERE a.session_id = ?;
`, sessionID)

	var username string
	var requestedForPromotion bool
	err = row.Scan(&username, &requestedForPromotion)
	if err == sql.ErrNoRows {
		// Session ID does not exist in the database
		info.User.LoggedIn = false
		return info
	} else if err != nil {
		// Error occurred while querying the database
		info.User.LoggedIn = false
		return info
	}

	// Check if the ID and email exists in the database
	row = db.QueryRow("SELECT uuid FROM users WHERE username = ?;", username)
	var uuid string
	err = row.Scan(&uuid)
	if err == sql.ErrNoRows {
		return info
	} else if err != nil {
		return info
	}

	info.User.ID = uuid
	info.User.Username = username
	info.User.PromotionRequest = requestedForPromotion
	info.User.LoggedIn = true
	info.User.TypeInt, info.User.Type, err = GetUserType(db, uuid)
	if err != nil {
		info.User.LoggedIn = false
		return info
	}
	info.User.UnreadNotifications, err = GetUnreadNotificationCount(db, uuid, info.User.TypeInt)
	if err != nil {
		info.User.LoggedIn = false
		return info
	}
	info.User.Admin.UnreadNotifications, err = GetAdminNotificationCount(db)
	if err != nil {
		info.User.LoggedIn = false
		return info
	}

	return info
}
