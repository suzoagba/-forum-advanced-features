package handlers

import (
	"database/sql"
	"fmt"
)

func GetUserType(db *sql.DB, uuid string) (int, string, error) {
	query := `
		SELECT level FROM users WHERE uuid = ?
	`
	row := db.QueryRow(query, uuid)

	var level int
	err := row.Scan(&level)
	if err != nil {
		return 0, "", err
	}

	switch level {
	case 0:
		return level, "User", nil
	case 1:
		return level, "Moderator", nil
	case 2:
		return level, "Administrator", nil
	default:
		return 0, "", fmt.Errorf("unknown user level: %d", level)
	}
}
