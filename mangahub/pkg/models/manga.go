package models

// Manga represents the core data structure for our system
type Manga struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"` // e.g., "Ongoing", "Completed"
	ChapterCount int      `json:"chapter_count"`
	Description  string   `json:"description"`
	Source       string   `json:"source"` // To track if it's Manual, API, or Scraped
}
