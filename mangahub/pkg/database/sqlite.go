package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
		author TEXT,
		genres TEXT,
		status TEXT,
		total_chapters INTEGER,
		chapter INTEGER,
		description TEXT
	);
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT,
		role TEXT
	);
	CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT,
    manga_id TEXT,
    current_chapter INTEGER,
    status TEXT,
    PRIMARY KEY(user_id, manga_id)
);`

	// Migration: Ensure `manga` has expected columns (for older DBs)
	mrows, err := db.Query("PRAGMA table_info(manga)")
	if err == nil {
		defer mrows.Close()
		hasGenres := false
		hasStatus := false
		hasTotal := false
		hasChapter := false
		hasDesc := false
		for mrows.Next() {
			var cid int
			var name string
			var ctype string
			var notnull int
			var dflt sql.NullString
			var pk int
			if err := mrows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
				continue
			}
			if name == "genres" {
				hasGenres = true
			}
			if name == "status" {
				hasStatus = true
			}
			if name == "total_chapters" {
				hasTotal = true
			}
			if name == "chapter" {
				hasChapter = true
			}
			if name == "description" {
				hasDesc = true
			}
		}

		if !hasGenres {
			_, _ = db.Exec("ALTER TABLE manga ADD COLUMN genres TEXT")
		}
		if !hasStatus {
			_, _ = db.Exec("ALTER TABLE manga ADD COLUMN status TEXT")
		}
		if !hasTotal {
			_, _ = db.Exec("ALTER TABLE manga ADD COLUMN total_chapters INTEGER")
		}
		if !hasChapter {
			_, _ = db.Exec("ALTER TABLE manga ADD COLUMN chapter INTEGER")
		}
		if !hasDesc {
			_, _ = db.Exec("ALTER TABLE manga ADD COLUMN description TEXT")
		}
	}

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Migration: Ensure `users` has `password_hash` column (fix from older schema)
	rows, err := db.Query("PRAGMA table_info(users)")
	if err == nil {
		defer rows.Close()
		hasPasswordHash := false
		hasPassword := false
		for rows.Next() {
			var cid int
			var name string
			var ctype string
			var notnull int
			var dflt sql.NullString
			var pk int
			if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
				continue
			}
			if name == "password_hash" {
				hasPasswordHash = true
			}
			if name == "password" {
				hasPassword = true
			}
		}

		if !hasPasswordHash {
			// Add the missing column
			if _, err := db.Exec("ALTER TABLE users ADD COLUMN password_hash TEXT"); err == nil {
				// If older table had `password` column, copy values over
				if hasPassword {
					_, _ = db.Exec("UPDATE users SET password_hash = password WHERE password_hash IS NULL OR password_hash = ''")
				}
			}
		}
	}

	// Migration: Normalize titles using 'Vol.' suffix to extract chapter numbers for older seeded data
	rows2, err := db.Query("SELECT id, title FROM manga WHERE title LIKE '%Vol.%'")
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id string
			var title string
			if err := rows2.Scan(&id, &title); err != nil {
				continue
			}
			parts := strings.Split(title, " Vol. ")
			if len(parts) == 2 {
				base := strings.TrimSpace(parts[0])
				numStr := strings.TrimSpace(parts[1])
				num, err2 := strconv.Atoi(numStr)
				if err2 != nil {
					continue
				}
				_, _ = db.Exec("UPDATE manga SET title = ?, chapter = ? WHERE id = ?", base, num, id)
			}
		}
	}

	return db, nil
}
