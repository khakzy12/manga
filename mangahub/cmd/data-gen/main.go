package main

import (
	"fmt"
	"log"
	"mangahub/pkg/utils" // Ensure this path matches your new structure
)

func main() {
	log.Println("🛠️ Generating 100 manual entries and fetching from MangaDex...")
	manualData := utils.GenerateManualEntries()
	apiData, _ := utils.FetchFromMangaDex()

	finalDb := append(manualData, apiData...)
	utils.SaveMangaToFile(finalDb, "data/manga.json") // Save to data folder

	fmt.Printf("✅ Success! Total entries in manga.json: %d\n", len(finalDb))
}
