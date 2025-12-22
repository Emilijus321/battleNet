package handlers

import (
	"battleNet/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"battleNet/external/tmdb"
	"battleNet/templates"

	"github.com/google/uuid"
)

// HandleSearchMovies - TMDB filmų paieška
func (h *Handler) HandleSearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if query == "" {
		result, err := h.tmdbClient.GetPopularMovies(r.Context(), page)
		if err != nil {
			log.Printf("Error getting popular movies: %v", err)
			http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
			return
		}

		email := h.sessionManager.GetString(r.Context(), "email")
		role := h.sessionManager.GetString(r.Context(), "role")

		// Tiesiogiai naudojame result.Results - tai jau []tmdb.TMDBMovie
		component := templates.SearchMoviesPage(email, role, result.Results, query, result.Page, result.TotalPages)
		component.Render(r.Context(), w)
		return
	}

	result, err := h.tmdbClient.SearchMovies(r.Context(), query, page)
	if err != nil {
		log.Printf("Error searching movies: %v", err)
		http.Error(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// Tiesiogiai naudojame result.Results
	component := templates.SearchMoviesPage(email, role, result.Results, query, result.Page, result.TotalPages)
	component.Render(r.Context(), w)
}

// HandleImportMovie - importuoti filmą iš TMDB į mūsų DB
func (h *Handler) HandleImportMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmdbIDStr := r.FormValue("tmdb_id")
	tmdbID, err := strconv.Atoi(tmdbIDStr)
	if err != nil {
		http.Error(w, "Invalid TMDB ID", http.StatusBadRequest)
		return
	}

	role := h.sessionManager.GetString(r.Context(), "role")
	if role != "admin" && role != "moderator" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmdbMovie, err := h.tmdbClient.GetMovieDetails(r.Context(), tmdbID)
	if err != nil {
		log.Printf("Error getting movie from TMDB: %v", err)
		http.Error(w, "Failed to fetch movie from TMDB", http.StatusInternalServerError)
		return
	}

	movie := convertToOurMovieModel(tmdbMovie)

	err = h.movieRepo.CreateMovie(r.Context(), &movie)
	if err != nil {
		log.Printf("Error importing movie: %v", err)
		http.Error(w, "Failed to import movie", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/movies", http.StatusSeeOther)
}

// HandleAPISearchMovies - API endpoint filmų paieškai
func (h *Handler) HandleAPISearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var result *tmdb.SearchResponse
	var err error

	if query == "" {
		result, err = h.tmdbClient.GetPopularMovies(r.Context(), page)
	} else {
		result, err = h.tmdbClient.SearchMovies(r.Context(), query, page)
	}

	if err != nil {
		log.Printf("Error in API search: %v", err)
		http.Error(w, `{"error": "Failed to search movies"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// convertToOurMovieModel - konvertuoja TMDB filmą į mūsų DB modelį
func convertToOurMovieModel(tmdbMovie *tmdb.TMDBMovie) models.Movie {
	movie := models.Movie{
		MovieID:      uuid.New(),
		Title:        tmdbMovie.Title,
		Overview:     stringPtr(tmdbMovie.Overview),
		ReleaseDate:  parseDate(tmdbMovie.ReleaseDate),
		PosterPath:   stringPtr(formatImageURL(tmdbMovie.PosterPath)),
		BackdropPath: stringPtr(formatImageURL(tmdbMovie.BackdropPath)),
		VoteAverage:  floatPtr(tmdbMovie.VoteAverage),
		VoteCount:    intPtr(tmdbMovie.VoteCount),
		Popularity:   floatPtr(tmdbMovie.Popularity),
		Runtime:      intPtr(tmdbMovie.Runtime),
		Status:       stringPtr(tmdbMovie.Status),
		ImdbID:       stringPtr(tmdbMovie.ImdbID),
		CreatedAt:    time.Now(),
	}
	return movie
}

// floatPtr - konvertuoja float64 į *float64
func floatPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

// intPtr - konvertuoja int į *int
func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// parseDate - konvertuoja string datą į *time.Time
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

// formatImageURL - sukuria pilną TMDB image URL
func formatImageURL(path string) string {
	if path == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + path
}

func extractYear(dateStr string) int {
	if len(dateStr) >= 4 {
		if year, err := strconv.Atoi(dateStr[:4]); err == nil {
			return year
		}
	}
	return 0
}

func formatPosterURL(path string) string {
	if path == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + path
}
