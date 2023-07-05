package handlers

import (
	"database/sql"
	"fmt"
)

func GetUserType(db *sql.DB, uuid string) (string, error) {
	query := `
		SELECT level FROM users WHERE uuid = ?
	`
	row := db.QueryRow(query, uuid)

	var level int
	err := row.Scan(&level)
	if err != nil {
		return "", err
	}

	switch level {
	case 0:
		return "User", nil
	case 1:
		return "Moderator", nil
	case 2:
		return "Administrator", nil
	default:
		return "", fmt.Errorf("unknown user level: %d", level)
	}
}
