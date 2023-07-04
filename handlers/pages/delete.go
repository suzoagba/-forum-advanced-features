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

		id := r.URL.Query().Get("id")
		editType := r.URL.Query().Get("type")

		var post structs.Post
		var comment structs.Comment
		if r.Method == http.MethodGet {

			log.Println("[edit] get")
			if editType == "post" {
				post, _ = GetPost(w, db, id)
			} else if editType == "comment" {
				comment, _ = GetComment(w, db, id)
			}

			// Create a data struct to pass to the template
			data := forPage{
				User:    user,
				Post:    post,
				Comment: comment,
				Tags:    database.Tags,
			}
			handlers.RenderTemplates("delete", data, w, r)

		} else if r.Method == http.MethodPost {
			//handlers.RenderTemplates("delete", forPage, w, r) // TODO
		}
	}
}
