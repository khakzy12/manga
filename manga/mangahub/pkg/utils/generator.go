package utils

import (
	"fmt"
	"mangahub/pkg/models"
)

func GenerateManualEntries() []models.Manga {
	var list []models.Manga

	// Pre-defined data to ensure genre requirements are met
	genres := []string{"Shounen", "Shoujo", "Seinen", "Josei"}
	statuses := []string{"Ongoing", "Completed", "Hiatus"}

	// Sample base titles to replicate (you can add more unique ones here)
	baseTitles := []string{
		"One Piece", "Naruto", "Bleach", "Nana", "Fruits Basket",
		"Monster", "Berserk", "Midnight Secretary", "Blue Box", "Vagabond",
	}

	fmt.Println("🛠️  Generating 100 manual entries...")

	for i := 1; i <= 100; i++ {
		// Rotate through genres and titles to create variety
		genre := genres[i%4]
		title := fmt.Sprintf("%s Vol. %d", baseTitles[i%10], (i/10)+1)

		m := models.Manga{
			ID:           i,
			Title:        title,
			Author:       fmt.Sprintf("Author %d", i),
			Genres:       []string{genre, "Action", "Drama"}, // Mixes the required genre with others
			Status:       statuses[i%3],
			ChapterCount: 10 + (i * 2),
			Description:  fmt.Sprintf("This is a manually entered description for %s.", title),
			Source:       "Manual",
		}
		list = append(list, m)
	}

	return list
}
