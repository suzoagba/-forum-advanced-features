package pages

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func NotifyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		if r.Method == http.MethodGet {
			forPage.Notifications, _ = handlers.GetUserNotifications(db, forPage.User.ID)
			err := handlers.MarkNotificationsAsRead(db, forPage.User.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			handlers.RenderTemplates("notify", forPage, w, r)
		}
	}
}
