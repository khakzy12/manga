package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite" // Use the pure-go driver
)

func InitDB() (*sql.DB, error) {
	// 1. Ensure the 'data' directory exists in the project root
	dbDir := "data"
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// 2. Define the path (e.g., data/mangahub.db)
	dbPath := filepath.Join(dbDir, "mangahub.db")

	// 3. Open using the 'sqlite' driver name (specific to glebarez)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 4. Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	// 5. Initialize Schema (Tables)
	query := `
	CREATE TABLE IF NOT EXISTS manga (
		id TEXT PRIMARY KEY,
		title TEXT,
		password_hash TEXT, 
		author TEXT
	);
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);
	CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT,
    manga_id TEXT,
    current_chapter INTEGER,
    status TEXT,
    PRIMARY KEY(user_id, manga_id)
);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}
