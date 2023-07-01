package database

import (
	"database/sql"
	"forum/structs"
	"strconv"
	"strings"
)

func GetAllPosts(db *sql.DB, tagID int, user structs.User) ([]structs.Post, error) {
	if tagID == 98 {
		return GetPostsCreatedByUser(db, user.Username)
	} else if tagID == 99 {
		return GetPostsLikedByUser(db, user.ID, 1)
	}

	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		GROUP BY p.postID
	`

	if tagID > 0 {
		query += `
			HAVING GROUP_CONCAT(t.id) LIKE '%` + strconv.Itoa(tagID) + `%'
		`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsCreatedByUser(db *sql.DB, username string) ([]structs.Post, error) {
	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		WHERE p.username = ?
		GROUP BY p.postID
	`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostsLikedByUser(db *sql.DB, userID string, reaction int) ([]structs.Post, error) {
	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		JOIN post_reactions pr ON p.postID = pr.post_id
		WHERE pr.user_id = ? AND pr.reaction_type = ?
		GROUP BY p.postID
	`

	rows, err := db.Query(query, userID, reaction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetCommentsByUser retrieves all comments made by a given user
func GetCommentsByUser(db *sql.DB, username string) ([]structs.CommentListing, error) {
	// Query to get all comments made by the given username
	query := `SELECT p.postID, p.username, p.title, p.description, p.creationDate,
                    c.commentID, c.content, c.postID, c.username, c.creationDate
              FROM posts p
              JOIN comments c ON p.postID = c.postID
              WHERE c.username = ?`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentListings := make([]structs.CommentListing, 0)
	postCommentsMap := make(map[int][]structs.Comment)
	postTitlesMap := make(map[int]string)

	for rows.Next() {
		var postID, commentID, commentPostID int
		var postUsername, title, description, postCreationDate string
		var commentContent, commentUsername, commentCreationDate string

		err := rows.Scan(&postID, &postUsername, &title, &description, &postCreationDate,
			&commentID, &commentContent, &commentPostID, &commentUsername, &commentCreationDate)
		if err != nil {
			return nil, err
		}

		comment := structs.Comment{
			ID:           strconv.Itoa(commentID),
			Content:      commentContent,
			PostID:       commentPostID,
			Username:     commentUsername,
			CreationDate: commentCreationDate,
		}

		postCommentsMap[postID] = append(postCommentsMap[postID], comment)
		postTitlesMap[postID] = title
	}

	for postID, comments := range postCommentsMap {
		post := structs.Post{
			ID:       postID,
			Username: comments[0].Username,
			Title:    postTitlesMap[postID],
		}
		commentListing := structs.CommentListing{
			Post:     post,
			Comments: comments,
		}
		commentListings = append(commentListings, commentListing)
	}

	return commentListings, nil
}

// GetCommentsByUserReaction retrieves comments reacted by the given username based on the reaction type
func GetCommentsByUserReaction(db *sql.DB, username string, reaction int) ([]structs.CommentListing, error) {
	// Query to get comments reacted by the given username and reaction type
	query := `SELECT p.postID, p.username, p.title, p.description, p.creationDate,
                    c.commentID, c.content, c.postID, c.username, c.creationDate
              FROM posts p
              JOIN comments c ON p.postID = c.postID
              JOIN comment_reactions cr ON c.commentID = cr.comment_id
              WHERE cr.user_id = ? AND cr.reaction_type = ?`

	rows, err := db.Query(query, username, reaction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentListings := make([]structs.CommentListing, 0)
	postCommentsMap := make(map[int][]structs.Comment)
	postTitlesMap := make(map[int]string)

	for rows.Next() {
		var postID, commentID, commentPostID int
		var postUsername, title, description, postCreationDate string
		var commentContent, commentUsername, commentCreationDate string

		err := rows.Scan(&postID, &postUsername, &title, &description, &postCreationDate,
			&commentID, &commentContent, &commentPostID, &commentUsername, &commentCreationDate)
		if err != nil {
			return nil, err
		}

		comment := structs.Comment{
			ID:           strconv.Itoa(commentID),
			Content:      commentContent,
			PostID:       commentPostID,
			Username:     commentUsername,
			CreationDate: commentCreationDate,
		}

		postCommentsMap[postID] = append(postCommentsMap[postID], comment)
		postTitlesMap[postID] = title
	}

	for postID, comments := range postCommentsMap {
		post := structs.Post{
			ID:       postID,
			Username: comments[0].Username,
			Title:    postTitlesMap[postID],
		}
		commentListing := structs.CommentListing{
			Post:     post,
			Comments: comments,
		}
		commentListings = append(commentListings, commentListing)
	}

	return commentListings, nil
}
