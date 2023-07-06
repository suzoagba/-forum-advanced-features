package user

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func UserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User

		if r.Method == http.MethodGet {
			name := r.URL.Query().Get("name")
			if forPage.User.LoggedIn {
				var err error
				forPage.UserInfo.Email, forPage.UserInfo.TypeInt, forPage.UserInfo.Type, err = GetUserInfo(db, name)
				if err != nil {
					handlers.ErrorHandler(w, http.StatusInternalServerError, "Error reading user info")
					return
				}
				forPage.UserInfo.Username = name
				handlers.RenderTemplates("user", forPage, w, r)
				return
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
			}
		} else if r.Method == http.MethodPost {
			if forPage.User.TypeInt == 2 {
				name := r.FormValue("name")
				level := r.FormValue("level")
				err := changeUserLevel(db, name, level)
				if err != nil {
					handlers.ErrorHandler(w, http.StatusInternalServerError, "Error changing user type")
					return
				}
				http.Redirect(w, r, "/user?name="+name, http.StatusFound)
				return
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
	}
}

func GetUserInfo(db *sql.DB, name string) (string, int, string, error) {
	query := `
		SELECT uuid, email, level FROM users WHERE username = ?
	`
	row := db.QueryRow(query, name)

	var uuid, email, uType string
	var level int
	err := row.Scan(&uuid, &email, &level)
	if err != nil {
		return "", 0, "", err
	}
	level, uType, err = handlers.GetUserType(db, uuid)

	return email, level, uType, nil
}

func changeUserLevel(db *sql.DB, name string, level string) error {
	updateQuery := `
		UPDATE users
		SET level = ?
		WHERE username = ?
	`
	_, err := db.Exec(updateQuery, level, name)
	if err != nil {
		return err
	}

	return nil
}
