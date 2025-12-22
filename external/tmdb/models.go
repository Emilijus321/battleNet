package tmdb

// Galite čia dėti papildomas TMDB struktūras
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
}

type CrewMember struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Job        string `json:"job"`
	Department string `json:"department"`
}
