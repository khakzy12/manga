package models

// ProgressUpdate must be Capitalized to be exported
type ProgressUpdate struct {
	Username string `json:"username"`
	MangaID  string `json:"manga_id"`
	Progress string `json:"progress"`
	Chapter  string `json:"chapter"`
}
