package handlers

import (
	"battleNet/external/tmdb"
	"battleNet/models"
	"battleNet/repository"
	"battleNet/templates"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

type Handler struct {
	userRepo       *repository.UserRepository
	movieRepo      *repository.MovieRepository
	reviewRepo     *repository.ReviewRepository
	watchlistRepo  *repository.WatchlistRepository
	jwtSecret      string
	sessionManager *scs.SessionManager
	tmdbClient     *tmdb.Client
}

func NewHandler(
	userRepo *repository.UserRepository,
	movieRepo *repository.MovieRepository,
	reviewRepo *repository.ReviewRepository,
	watchlistRepo *repository.WatchlistRepository,
	jwtSecret string,
	sessionManager *scs.SessionManager,
	tmdbClient *tmdb.Client,
) *Handler {
	return &Handler{
		userRepo:       userRepo,
		movieRepo:      movieRepo,
		reviewRepo:     reviewRepo,
		watchlistRepo:  watchlistRepo,
		jwtSecret:      jwtSecret,
		sessionManager: sessionManager,
		tmdbClient:     tmdbClient,
	}
}

// HandleAdminMovies - admin filmų valdymo puslapis
func (h *Handler) HandleAdminMovies(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// Log debug info
	log.Printf("AdminMovies accessed by: %s (role: %s)", email, role)

	// Gauti visus filmus
	movies, err := h.movieRepo.GetMovies(r.Context(), 100, 0)
	if err != nil {
		log.Printf("Error getting movies: %v", err)
		// Grąžinti tuščią sąrašą
		movies = []models.Movie{}
	}

	log.Printf("Found %d movies", len(movies))

	// Renderinti template
	component := templates.AdminMoviesPage(email, role, movies)
	component.Render(r.Context(), w)
}

// ==================== FILMO KŪRIMAS ====================

func (h *Handler) HandleCreateMoviePage(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// NAUDOJAME TEMPLATE'Ą iš create_movie.templ
	component := templates.CreateMoviePage(email, role)
	component.Render(r.Context(), w)
}

// HandleCreateMovie - apdoroja filmo kūrimą
func (h *Handler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Gauti duomenis (BE poster_path ir backdrop_path)
	title := r.FormValue("title")
	overview := r.FormValue("overview")
	releaseDateStr := r.FormValue("release_date")
	voteAvgStr := r.FormValue("vote_average")
	voteCountStr := r.FormValue("vote_count")
	popularityStr := r.FormValue("popularity")
	runtimeStr := r.FormValue("runtime")
	status := r.FormValue("status")
	imdbID := r.FormValue("imdb_id") // Tik IMDB ID

	// Konvertuoti
	var releaseDate *time.Time
	if releaseDateStr != "" {
		rd, _ := time.Parse("2006-01-02", releaseDateStr)
		releaseDate = &rd
	}

	var voteAvg *float64
	if voteAvgStr != "" {
		va, _ := strconv.ParseFloat(voteAvgStr, 64)
		voteAvg = &va
	}

	var voteCount *int
	if voteCountStr != "" {
		vc, _ := strconv.Atoi(voteCountStr)
		voteCount = &vc
	} else {
		defaultVoteCount := 0
		voteCount = &defaultVoteCount
	}

	var popularity *float64
	if popularityStr != "" {
		pop, _ := strconv.ParseFloat(popularityStr, 64)
		popularity = &pop
	} else {
		defaultPopularity := 0.0
		popularity = &defaultPopularity
	}

	var runtime *int
	if runtimeStr != "" {
		rt, _ := strconv.Atoi(runtimeStr)
		runtime = &rt
	} else {
		defaultRuntime := 120
		runtime = &defaultRuntime
	}

	if status == "" {
		status = "Released"
	}

	// Sukurti filmą (poster_path ir backdrop_path = nil)
	movie := &models.Movie{
		Title:        title,
		Overview:     stringPtr(overview),
		ReleaseDate:  releaseDate,
		PosterPath:   nil, // Nėra nuotraukų
		BackdropPath: nil, // Nėra nuotraukų
		VoteAverage:  voteAvg,
		VoteCount:    voteCount,
		Popularity:   popularity,
		Runtime:      runtime,
		Status:       &status,
		ImdbID:       stringPtr(imdbID), // Tik IMDB ID
	}

	err := h.movieRepo.CreateMovie(r.Context(), movie)
	if err != nil {
		log.Printf("Error creating movie: %v", err)
		http.Error(w, "Failed to create movie", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/movies", http.StatusSeeOther)
}

// ==================== FILMO REDAGAVIMAS ====================

// HandleEditMoviePage - filmo redagavimo forma
func (h *Handler) HandleEditMoviePage(w http.ResponseWriter, r *http.Request) {
	// Gauti filmo ID iš URL
	movieIDStr := r.URL.Query().Get("id")
	if movieIDStr == "" {
		http.Error(w, "Movie ID required", http.StatusBadRequest)
		return
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Gauti filmą iš DB
	movie, err := h.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// NAUDOTI TEMPLATE'Ą (ne hardcoded HTML)
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	component := templates.EditMoviePage(email, role, *movie)
	component.Render(r.Context(), w)
}

// HandleUpdateMovie - apdoroja filmo atnaujinimą
func (h *Handler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Gauti filmo ID
	movieIDStr := r.FormValue("movie_id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Gauti duomenis
	title := r.FormValue("title")
	overview := r.FormValue("overview")
	releaseDateStr := r.FormValue("release_date")
	voteAvgStr := r.FormValue("vote_average")
	runtimeStr := r.FormValue("runtime")
	status := r.FormValue("status")

	// Konvertuoti
	var releaseDate *time.Time
	if releaseDateStr != "" {
		rd, _ := time.Parse("2006-01-02", releaseDateStr)
		releaseDate = &rd
	}

	var voteAvg *float64
	if voteAvgStr != "" {
		va, _ := strconv.ParseFloat(voteAvgStr, 64)
		voteAvg = &va
	}

	var runtime *int
	if runtimeStr != "" {
		rt, _ := strconv.Atoi(runtimeStr)
		runtime = &rt
	}

	if status == "" {
		status = "Released"
	}

	// Atnaujinti filmą
	movie := &models.Movie{
		Title:       title,
		Overview:    stringPtr(overview),
		ReleaseDate: releaseDate,
		VoteAverage: voteAvg,
		Runtime:     runtime,
		Status:      &status,
	}

	err = h.movieRepo.UpdateMovie(r.Context(), movieID, movie)
	if err != nil {
		log.Printf("Error updating movie: %v", err)
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/movies", http.StatusSeeOther)
}

// ==================== FILMO IŠTRYNIMAS ====================

// HandleDeleteMovie - ištrina filmą
func (h *Handler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	movieIDStr := r.FormValue("movie_id")
	if movieIDStr == "" {
		http.Error(w, "Movie ID required", http.StatusBadRequest)
		return
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Ištrinti filmą
	err = h.movieRepo.DeleteMovie(r.Context(), movieID)
	if err != nil {
		log.Printf("Error deleting movie: %v", err)
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/movies", http.StatusSeeOther)
}

// ==================== API ENDPOINTS ====================

func (h *Handler) HandleAPICreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"API endpoint for creating movies"}`))
}

func (h *Handler) HandleAPIUpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"Movie updated"}`))
}

func (h *Handler) HandleAPIDeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"Movie deleted"}`))
}

// ==================== MODERATOR HANDLERS ====================

// HandleModeratorDashboard - moderator dashboard
func (h *Handler) HandleModeratorDashboard(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")
	name := h.sessionManager.GetString(r.Context(), "name")
	username := h.sessionManager.GetString(r.Context(), "username")

	component := templates.DashboardPage(email, "", role, name, username)
	component.Render(r.Context(), w)
}

// HandleModeratorUsers - show all users for moderator
func (h *Handler) HandleModeratorUsers(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")
	currentUserIDStr := h.sessionManager.GetString(r.Context(), "userID")

	// Get all users
	users, err := h.userRepo.GetAllUsers(r.Context(), 50, 0)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		users = []models.User{}
	}

	component := templates.ModeratorUsersPage(email, role, users, currentUserIDStr)
	component.Render(r.Context(), w)
}

// HandleModeratorUpdateRole - update user role
func (h *Handler) HandleModeratorUpdateRole(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("user_id")
	newRole := r.FormValue("new_role")
	currentUserIDStr := h.sessionManager.GetString(r.Context(), "userID")

	// Check if trying to change own role
	if userIDStr == currentUserIDStr {
		http.Error(w, "Cannot change your own role", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Validate role
	validRoles := map[string]bool{"user": true, "moderator": true, "admin": true}
	if !validRoles[newRole] {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Update role
	err = h.userRepo.UpdateUserRole(r.Context(), userID, newRole)
	if err != nil {
		log.Printf("Error updating user role: %v", err)
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/moderator/users", http.StatusSeeOther)
}

// HandleModeratorDeactivateUser - deactivate user
func (h *Handler) HandleModeratorDeactivateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("user_id")
	currentUserIDStr := h.sessionManager.GetString(r.Context(), "userID")

	// Check if trying to deactivate self
	if userIDStr == currentUserIDStr {
		http.Error(w, "Cannot deactivate yourself", http.StatusForbidden)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Deactivate user
	err = h.userRepo.DeactivateUser(r.Context(), userID)
	if err != nil {
		log.Printf("Error deactivating user: %v", err)
		http.Error(w, "Failed to deactivate user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/moderator/users", http.StatusSeeOther)
}

// ==================== MODERATOR API ENDPOINTS ====================

func (h *Handler) HandleAPIModeratorUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"API endpoint for moderator users"}`))
}

func (h *Handler) HandleAPIModeratorUpdateRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"User role updated"}`))
}

func (h *Handler) HandleAPIModeratorDeactivateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success", "message":"User deactivated"}`))
}

// ==================== HELPER FUNCTIONS ====================

// Helper function for string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper for selected attribute
func selectedIf(condition bool, value string) string {
	if condition {
		return value
	}
	return ""
}
