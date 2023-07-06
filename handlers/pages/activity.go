package pages

import (
	"database/sql"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"net/http"
)

func ActivityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		forPage.Activity.Posts, _ = database.GetPostsCreatedByUser(db, forPage.User.Username)
		forPage.Activity.Comments, _ = database.GetCommentsByUser(db, forPage.User.Username)
		forPage.Activity.PostLikes, _ = database.GetPostsLikedByUser(db, forPage.User.ID, 1)
		forPage.Activity.PostDislikes, _ = database.GetPostsLikedByUser(db, forPage.User.ID, 0)
		forPage.Activity.CommentLikes, _ = database.GetCommentsByUserReaction(db, forPage.User.ID, 1)
		forPage.Activity.CommentDislikes, _ = database.GetCommentsByUserReaction(db, forPage.User.ID, 0)
		if r.Method == http.MethodGet {
			handlers.RenderTemplates("activity", forPage, w, r)
		}
	}
}
