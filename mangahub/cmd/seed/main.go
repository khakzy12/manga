// cmd/seed/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mangahub/pkg/database"
	"mangahub/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🌱 Starting database seeding process...")

	// This calls the InitDB we updated above, which points to data/mangahub.db
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Insert generated manga entries (100)
	entries := utils.GenerateManualEntries()
	inserted := 0
	for _, m := range entries {
		genresJson, _ := json.Marshal(m.Genres)
		idStr := fmt.Sprintf("%d", m.ID)
		_, err := db.Exec("INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, chapter, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			idStr, m.Title, m.Author, string(genresJson), m.Status, m.ChapterCount, m.Chapter, m.Description)
		if err != nil {
			log.Printf("⚠️ Error inserting manga %s: %v", m.Title, err)
			continue
		}
		inserted++
	}
	fmt.Printf("✅ Successfully seeded %d manga entries into data/mangahub.db\n", inserted)

	// Helper to seed users with bcrypt hashed passwords
	seedUser := func(username, password, role string) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("⚠️ Could not hash password for %s: %v", username, err)
			return
		}
		id := uuid.NewString()
		_, err = db.Exec("INSERT OR IGNORE INTO users (id, username, password_hash, role) VALUES (?, ?, ?, ?)",
			id, username, string(hashed), role)
		if err != nil {
			log.Printf("⚠️ Error inserting user %s: %v", username, err)
			return
		}
		fmt.Printf("✅ Seeded user: %s (role=%s) id=%s\n", username, role, id)
	}

	// Seed an admin and a demo user
	seedUser("admin", "adminpass", "admin")
	seedUser("demo", "demopass", "user")

	fmt.Println("✅ Seeding complete")
}
