package database

import "fmt"

// Statetement to create tables using sqlite
var TableCreationStatements = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		
		nickname TEXT UNIQUE NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);`,
	`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (post_id) REFERENCES posts (id)
	);`,
	`CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		is_like INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		UNIQUE(user_id, post_id)
	);`,
	`CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);`,
	`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`,
	`CREATE TABLE IF NOT EXISTS comment_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		is_like INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		UNIQUE(user_id, comment_id)
	);`,
	`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		message TEXT DEFAULT "",
		is_read BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT (DATETIME('now', 'localtime')),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`,
	`INSERT OR IGNORE INTO categories (name) VALUES 
		('Technology'),
		('Sports'),
		('Entertainment'),
		('Science'),
		('Politics'),
		('Health'),
		('Travel'),
		('Food'),
		('Other');`,
}

// Create tables in the database
func CreateTables(statements []string) (err error) {
	for _, statement := range statements {
		_, err = DB.Exec(statement)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}
	return
}
