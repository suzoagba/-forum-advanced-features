package pages

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func EditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		if r.Method == http.MethodGet {
			handlers.RenderTemplates("edit", forPage, w, r)
		} else if r.Method == http.MethodPost {
			handlers.RenderTemplates("edit", forPage, w, r) // TODO
		}
	}
}
