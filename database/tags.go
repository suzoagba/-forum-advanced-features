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
