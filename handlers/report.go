package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func ReportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		user := IsLoggedIn(r, db).User
		if r.Method != http.MethodPost || user.TypeInt != 1 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			id := r.FormValue("id")
			approved := r.FormValue("approved")
			if approved == "true" {
				fmt.Println(id)
			} else {
				reason := r.FormValue("reason")
				fmt.Println(reason)
			}
		}
	}
}
