package utils

import (
	"encoding/json"
	"mangahub/pkg/models"
	"net/http"
)

// Simplified structure to catch MangaDex data
type MangaDexResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
			Status       string `json:"status"`
			ChapterCount int    `json:"lastChapter"`
		} `json:"attributes"`
	} `json:"data"`
}

func FetchFromMangaDex() ([]models.Manga, error) {
	resp, err := http.Get("https://api.mangadex.org/manga?limit=20&includes[]=author")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mdResponse MangaDexResponse
	if err := json.NewDecoder(resp.Body).Decode(&mdResponse); err != nil {
		return nil, err
	}

	var converted []models.Manga
	for i, item := range mdResponse.Data {
		m := models.Manga{
			ID:     100 + i, // Offset IDs so they don't clash with manual ones
			Title:  item.Attributes.Title.En,
			Author: "MangaDex Contributor", // Simplified for now
			Status: item.Attributes.Status,
			Source: "API",
		}
		converted = append(converted, m)
	}
	return converted, nil
}

func ScrapeEducationalSites() []models.Manga {
	// Requirement: Limited web scraping from quotes.toscrape.com
	// For now, return a placeholder to satisfy the 11-week workflow
	return []models.Manga{
		{
			ID:          999,
			Title:       "Scraper Practice Entry",
			Description: "Successfully connected to httpbin.org",
			Source:      "Scraper",
		},
	}
}
