package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"battleNet/models"
	"battleNet/templates"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandleHome displays home page
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	component := templates.HomePage()
	component.Render(r.Context(), w)
}

// HandleDashboard displays user dashboard
func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	userID := h.sessionManager.GetString(r.Context(), "userID")
	role := h.sessionManager.GetString(r.Context(), "role")
	name := h.sessionManager.GetString(r.Context(), "name")

	component := templates.DashboardPage(email, userID, role, name)
	component.Render(r.Context(), w)
}

// HandleProfile displays user profile
func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	name := h.sessionManager.GetString(r.Context(), "name")
	role := h.sessionManager.GetString(r.Context(), "role")

	component := templates.ProfilePage(email, name, role)
	component.Render(r.Context(), w)
}

// HandleMovies displays movies list
func (h *Handler) HandleMovies(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get movies from database
	movies, err := h.movieRepo.GetMovies(r.Context(), int32(limit), int32(offset))
	if err != nil {
		log.Printf("Error getting movies: %v", err)
		http.Error(w, "Failed to load movies", http.StatusInternalServerError)
		return
	}

	component := templates.MoviesPage(email, role, movies)
	component.Render(r.Context(), w)
}

// HandleMovieDetail displays movie details
func (h *Handler) HandleMovieDetail(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	movieIDStr := chi.URLParam(r, "id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Get movie from database
	movie, err := h.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Get reviews for this movie
	reviews, err := h.reviewRepo.GetMovieReviews(r.Context(), movieID)
	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		// Continue without reviews
		reviews = []models.Review{}
	}

	// Check if movie is in user's watchlist
	var inWatchlist bool
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	if userIDStr != "" {
		userID, _ := uuid.Parse(userIDStr)
		inWatchlist, _ = h.watchlistRepo.CheckWatchlist(r.Context(), userID, movieID)
	}

	component := templates.MovieDetailPage(email, role, *movie, reviews, inWatchlist)
	component.Render(r.Context(), w)
}

// HandleAPIMovies returns movies as JSON (API endpoint)
func (h *Handler) HandleAPIMovies(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	movies, err := h.movieRepo.GetMovies(r.Context(), int32(limit), int32(offset))
	if err != nil {
		log.Printf("Error getting movies for API: %v", err)
		http.Error(w, `{"error": "Failed to fetch movies"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"movies": movies,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(movies), // Note: You might want to add a count method
		},
	})
}

// HandleAPIMovieDetail returns movie details as JSON (API endpoint)
func (h *Handler) HandleAPIMovieDetail(w http.ResponseWriter, r *http.Request) {
	movieIDStr := chi.URLParam(r, "id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid movie ID"}`, http.StatusBadRequest)
		return
	}

	movie, err := h.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		http.Error(w, `{"error": "Movie not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}
