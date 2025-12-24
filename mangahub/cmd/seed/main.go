// cmd/seed/main.go
package main

import (
	"fmt"
	"log"
	"mangahub/pkg/database"
)

func main() {
	fmt.Println("🌱 Starting database seeding process...")

	// This calls the InitDB we updated above, which points to data/mangahub.db
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Example: Insert a test manga
	_, err = db.Exec("INSERT OR IGNORE INTO manga (id, title, author) VALUES (?, ?, ?)",
		"1", "One Piece", "Eiichiro Oda")

	if err != nil {
		log.Printf("⚠️ Error inserting data: %v", err)
	} else {
		fmt.Println("✅ Successfully seeded data into data/mangahub.db")
	}
}
