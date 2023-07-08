package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func StartDB(input bool) {
	if input {
		err := createDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createDatabase() error {
	// Create the database folder if it doesn't exist
	err := os.MkdirAll("./database", os.ModePerm)
	if err != nil {
		return err
	}

	// Create the database file
	_, err = os.Create("./database/forum.db")
	if err != nil {
		return err
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
		uuid TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		username TEXT UNIQUE,
		password TEXT,
		level INTEGER DEFAULT 0,
		requested_for_promotion BOOLEAN DEFAULT false
	);

		CREATE TABLE IF NOT EXISTS posts (
			postID INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			title TEXT,
			description TEXT,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			imageFilename TEXT,
			edited BOOLEAN DEFAULT false,
			timeEdited TIMESTAMP,
			approved BOOLEAN DEFAULT true,
			reported BOOLEAN DEFAULT false,
			report_reason TEXT,
			FOREIGN KEY (username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS comments (
			commentID INTEGER PRIMARY KEY AUTOINCREMENT,
			postID INTEGER,
			username TEXT,
			content TEXT,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			edited BOOLEAN DEFAULT false,
			timeEdited TIMESTAMP,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS authenticated_users (
			session_id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			FOREIGN KEY(username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);

		CREATE TABLE IF NOT EXISTS post_tags (
			postID INTEGER,
			tagID INTEGER,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (tagID) REFERENCES tags(id),
			PRIMARY KEY (postID, tagID)
		);

		INSERT OR IGNORE INTO tags (name) VALUES ('Cooking'), ('Mechanics'), ('Travel'), ('IT'), ('Random'), ('Market');

		CREATE TABLE IF NOT EXISTS post_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id UUID,
			reaction_type BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id),
			FOREIGN KEY (post_id) REFERENCES posts(postID),
			FOREIGN KEY (user_id) REFERENCES users(uuid)
		);
		
		CREATE TABLE IF NOT EXISTS comment_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			comment_id INTEGER,
			user_id UUID,
			reaction_type BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(comment_id, user_id),
			FOREIGN KEY (comment_id) REFERENCES comments(commentID),
			FOREIGN KEY (user_id) REFERENCES users(uuid)
		);	

		CREATE TABLE IF NOT EXISTS notifications (
			  notificationID INTEGER PRIMARY KEY AUTOINCREMENT,
			  userID UUID,
			  whoID UUID,
			  actionDone TEXT,
			  isPost BOOLEAN DEFAULT false,
			  isComment BOOLEAN DEFAULT false,
			  postID INTEGER,
			  commentID INTEGER,
			  isRead BOOLEAN DEFAULT false,
			  createdDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			  FOREIGN KEY (userID) REFERENCES users(uuid),
			  FOREIGN KEY (whoID) REFERENCES users(uuid),
			  FOREIGN KEY (postID) REFERENCES posts(postID),
			  FOREIGN KEY (commentID) REFERENCES comments(commentID)
		); 

		CREATE TABLE IF NOT EXISTS admin_notifications (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    post BOOLEAN DEFAULT false,
		    postID INTEGER DEFAULT 0,
		    user BOOLEAN DEFAULT false,
		    userID UUID,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (userID) REFERENCES users(uuid)
		);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	fmt.Println("Database and tables created successfully!")
	return nil
}
