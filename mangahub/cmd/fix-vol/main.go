package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	// Open a dedicated connection with a busy timeout to tolerate transient locks
	dbPath := "data/mangahub.db?_busy_timeout=5000"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure PRAGMA busy_timeout is set as well
	_, _ = db.Exec("PRAGMA busy_timeout = 5000")

	// 1) Read matching rows into memory (avoid holding a read cursor while updating)
	rows, err := db.Query("SELECT id, title FROM manga WHERE title LIKE '%Vol.%'")
	if err != nil {
		log.Fatal(err)
	}
	var targets []struct{ id, title string }
	for rows.Next() {
		var id, title string
		if err := rows.Scan(&id, &title); err != nil {
			continue
		}
		targets = append(targets, struct{ id, title string }{id, title})
	}
	rows.Close()

	if len(targets) == 0 {
		fmt.Println("No 'Vol.' titles found to normalize.")
		return
	}

	// 2) Update each row with retry logic
	for _, t := range targets {
		fmt.Printf("Found: %s -> %s\n", t.id, t.title)
		parts := strings.Split(t.title, " Vol. ")
		if len(parts) != 2 {
			continue
		}
		base := strings.TrimSpace(parts[0])
		numStr := strings.TrimSpace(parts[1])
		num, e := strconv.Atoi(numStr)
		if e != nil {
			fmt.Println("parse error", e)
			continue
		}

		// Retry loop for SQLITE_BUSY
		var lastErr error
		for attempt := 0; attempt < 8; attempt++ {
			_, err := db.Exec("UPDATE manga SET title = ?, chapter = ? WHERE id = ?", base, num, t.id)
			if err == nil {
				fmt.Printf("updated id=%s -> title=%q chapter=%d\n", t.id, base, num)
				lastErr = nil
				break
			}
			lastErr = err
			if strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "SQLITE_BUSY") {
				time.Sleep(time.Duration(200*(attempt+1)) * time.Millisecond)
				continue
			}
			break
		}
		if lastErr != nil {
			fmt.Printf("failed to update id=%s: %v\n", t.id, lastErr)
		}
	}

	fmt.Println("Normalization complete")
}
