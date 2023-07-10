package database

import (
	"database/sql"
	"forum/structs"
	"log"
)

var Tags []structs.Tag

func GetTags(db *sql.DB) {
	tags, err := getTagsFromDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
	Tags = tags
}

func getTagsFromDatabase(db *sql.DB) ([]structs.Tag, error) {
	query := "SELECT id, name FROM tags"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []structs.Tag
	for rows.Next() {
		var tag structs.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func GetPostTags(db *sql.DB, postID string) ([]string, error) {
	query := `
		SELECT tagID FROM post_tags WHERE postID = $1;
	`

	// Execute the query
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result rows and retrieve the tagIDs
	var tags []string
	for rows.Next() {
		var tagID string
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tagID)
	}

	// Check for any errors during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tags, nil
}

func AddTag(db *sql.DB, tagName string) error {
	// Insert the new tag into the tags table
	insertQuery := `
		INSERT INTO tags (name) VALUES (?)
	`
	_, err := db.Exec(insertQuery, tagName)
	if err != nil {
		return err
	}

	GetTags(db)
	return nil
}

func DeleteTag(db *sql.DB, tagID string) ([]string, error) {
	// Check if there are any posts associated only with the given tag
	checkPostsQuery := `
		SELECT postID
		FROM post_tags
		WHERE tagID = ? AND postID NOT IN (
			SELECT postID
			FROM post_tags
			WHERE tagID <> ?
		)
	`
	rows, err := db.Query(checkPostsQuery, tagID, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []string
	for rows.Next() {
		var postID string
		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	// If there are posts associated only with the given tag, return the list of postIDs
	if len(postIDs) > 0 {
		return postIDs, nil
	}

	// Delete the tag from the tags table
	deleteTagQuery := `
		DELETE FROM tags
		WHERE id = ?
	`
	_, err = db.Exec(deleteTagQuery, tagID)
	if err != nil {
		return nil, err
	}

	GetTags(db)
	return nil, nil
}

func GetTagNameByID(db *sql.DB, tagID string) (string, error) {
	query := `
		SELECT name FROM tags WHERE id = ?
	`
	row := db.QueryRow(query, tagID)

	var tagName string
	err := row.Scan(&tagName)
	if err != nil {
		return "", err
	}

	return tagName, nil
}
