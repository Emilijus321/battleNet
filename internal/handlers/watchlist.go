package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"battleNet/models"
	"battleNet/templates"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandleWatchlist displays user's watchlist
func (h *Handler) HandleWatchlist(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusUnauthorized)
		return
	}

	// Get watchlist from database
	watchlist, err := h.watchlistRepo.GetUserWatchlist(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting watchlist: %v", err)
		http.Error(w, "Failed to load watchlist", http.StatusInternalServerError)
		return
	}

	component := templates.WatchlistPage(email, role, watchlist)
	component.Render(r.Context(), w)
}

// HandleAddToWatchlist adds movie to watchlist
func (h *Handler) HandleAddToWatchlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	movieIDStr := r.FormValue("movie_id")

	if userIDStr == "" || movieIDStr == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusUnauthorized)
		return
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Check if movie exists
	_, err = h.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Add to watchlist
	_, err = h.watchlistRepo.AddToWatchlist(r.Context(), userID, movieID)
	if err != nil {
		log.Printf("Error adding to watchlist: %v", err)
		http.Error(w, "Failed to add to watchlist", http.StatusInternalServerError)
		return
	}

	// Redirect back to previous page or movies page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/movies"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// HandleAPIWatchlist returns watchlist as JSON (API endpoint)
func (h *Handler) HandleAPIWatchlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	watchlist, err := h.watchlistRepo.GetUserWatchlist(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting watchlist for API: %v", err)
		http.Error(w, `{"error": "Failed to fetch watchlist"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(watchlist)
}

// HandleAPIAddToWatchlist adds movie to watchlist (API endpoint)
func (h *Handler) HandleAPIAddToWatchlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var request struct {
		MovieID string `json:"movie_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	movieID, err := uuid.Parse(request.MovieID)
	if err != nil {
		http.Error(w, `{"error": "Invalid movie ID"}`, http.StatusBadRequest)
		return
	}

	item, err := h.watchlistRepo.AddToWatchlist(r.Context(), userID, movieID)
	if err != nil {
		log.Printf("Error adding to watchlist via API: %v", err)
		http.Error(w, `{"error": "Failed to add to watchlist"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// HandleAPIRemoveFromWatchlist removes movie from watchlist (API endpoint)
func (h *Handler) HandleAPIRemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	movieIDStr := chi.URLParam(r, "movieId")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid movie ID"}`, http.StatusBadRequest)
		return
	}

	err = h.watchlistRepo.RemoveFromWatchlist(r.Context(), userID, movieID)
	if err != nil {
		log.Printf("Error removing from watchlist via API: %v", err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<div hx-swap-oob="beforeend:#toast-container">
    <div class="toast toast-error">
        Failed to remove from watchlist
        <button class="close" onclick="this.parentElement.remove()">×</button>
    </div>
</div>
`))
		return
	}

	// Return both: the removed movie card and success toast
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<div hx-swap-oob="delete:.movie-card[data-movie-id='` + movieIDStr + `']"></div>
<div hx-swap-oob="beforeend:#toast-container">
    <div class="toast toast-success">
        Movie was successfully removed from watchlist
        <button class="close" onclick="this.parentElement.remove()">×</button>
    </div>
</div>
`))
}

// HandleRemoveFromWatchlist removes movie from watchlist (HTML form)
func (h *Handler) HandleRemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	movieIDStr := r.FormValue("movie_id")

	if userIDStr == "" || movieIDStr == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusUnauthorized)
		return
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = h.watchlistRepo.RemoveFromWatchlist(r.Context(), userID, movieID)
	if err != nil {
		log.Printf("Error removing from watchlist: %v", err)
		// Galite grąžinti error puslapį arba redirect su flash message
		http.Error(w, "Failed to remove from watchlist", http.StatusInternalServerError)
		return
	}

	// Redirect back to movie page
	http.Redirect(w, r, fmt.Sprintf("/movies/%s", movieIDStr), http.StatusSeeOther)
}
