// cmd/seed/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mangahub/pkg/database"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type MangaJSON struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	ChapterCount int      `json:"chapter_count"`
	Description  string   `json:"description"`
	Source       string   `json:"source"`
}

func main() {
	fmt.Println("🌱 Starting database seeding process...")

	// This calls the InitDB we updated above, which points to data/mangahub.db
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Read manga data from manga.json
	jsonFile, err := os.ReadFile("data/manga.json")
	if err != nil {
		log.Fatalf("❌ Failed to read manga.json: %v", err)
	}

	var mangaList []MangaJSON
	if err := json.Unmarshal(jsonFile, &mangaList); err != nil {
		log.Fatalf("❌ Failed to parse manga.json: %v", err)
	}

	fmt.Printf("📚 Found %d manga entries in manga.json\n", len(mangaList))

	// Insert all manga from JSON
	successCount := 0
	for _, manga := range mangaList {
		// Convert genres array to JSON string
		genresJSON, _ := json.Marshal(manga.Genres)

		_, err = db.Exec(`INSERT OR IGNORE INTO manga 
			(id, title, author, genres, status, total_chapters, description) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprint(manga.ID), manga.Title, manga.Author, string(genresJSON),
			manga.Status, manga.ChapterCount, manga.Description)

		if err != nil {
			log.Printf("⚠️ Error inserting manga %s: %v", manga.Title, err)
		} else {
			successCount++
		}
	}
	fmt.Printf("✅ Successfully seeded %d/%d manga entries\n", successCount, len(mangaList))

	// Create default admin user
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("⚠️ Error hashing password: %v", err)
		return
	}

	_, err = db.Exec("INSERT OR IGNORE INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		"admin", string(hashedPassword), "admin")

	if err != nil {
		log.Printf("⚠️ Error inserting admin user: %v", err)
	} else {
		fmt.Println("✅ Successfully created admin user")
		fmt.Println("   Username: admin")
		fmt.Println("   Password: admin123")
	}

	// Create default regular user
	userPassword := "user123"
	hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("⚠️ Error hashing user password: %v", err)
		return
	}

	_, err = db.Exec("INSERT OR IGNORE INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		"user", string(hashedUserPassword), "user")

	if err != nil {
		log.Printf("⚠️ Error inserting regular user: %v", err)
	} else {
		fmt.Println("✅ Successfully created regular user")
		fmt.Println("   Username: user")
		fmt.Println("   Password: user123")
	}
}
