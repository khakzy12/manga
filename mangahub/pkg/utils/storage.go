package utils

import (
	"encoding/json"
	"os"

	"mangahub/pkg/models" // Adjust "mangahub" to your module name
)

func SaveMangaToFile(mangaList []models.Manga, filename string) error {
	data, err := json.MarshalIndent(mangaList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadMangaFromFile(filename string) ([]models.Manga, error) {
	var mangaList []models.Manga
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &mangaList)
	return mangaList, err
}
