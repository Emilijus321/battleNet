package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"battleNet/models"
	//"battleNet/templates"

	"github.com/google/uuid"
)

// HandleCreateReview creates a new review
func (h *Handler) HandleCreateReview(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user session", http.StatusUnauthorized)
		return
	}

	movieID, err := uuid.Parse(r.FormValue("movie_id"))
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil || rating < 1 || rating > 10 {
		http.Error(w, "Rating must be between 1 and 10", http.StatusBadRequest)
		return
	}

	params := models.CreateReviewParams{
		UserID:           userID,
		MovieID:          movieID,
		Rating:           rating,
		Title:            r.FormValue("title"),
		Content:          r.FormValue("content"),
		ContainsSpoilers: r.FormValue("contains_spoilers") == "on",
		IsPublic:         true,
	}

	_, err = h.reviewRepo.CreateReview(r.Context(), params)
	if err != nil {
		log.Printf("Error creating review: %v", err)
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	// Redirect back to movie page
	http.Redirect(w, r, "/movies/"+movieID.String(), http.StatusSeeOther)
}

// HandleAPIReviews returns reviews as JSON (API endpoint)
func (h *Handler) HandleAPIReviews(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movie_id")

	var reviews []models.Review
	var err error

	if movieIDStr != "" {
		movieID, err := uuid.Parse(movieIDStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid movie ID"}`, http.StatusBadRequest)
			return
		}
		reviews, err = h.reviewRepo.GetMovieReviews(r.Context(), movieID)
	} else {
		// Get all public reviews or implement pagination
		http.Error(w, `{"error": "movie_id parameter is required"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error getting reviews for API: %v", err)
		http.Error(w, `{"error": "Failed to fetch reviews"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// HandleAPICreateReview creates a review (API endpoint)
func (h *Handler) HandleAPICreateReview(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var request struct {
		MovieID string `json:"movie_id"`
		Rating  int    `json:"rating"`
		Title   string `json:"title"`
		Content string `json:"content"`
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

	if request.Rating < 1 || request.Rating > 10 {
		http.Error(w, `{"error": "Rating must be between 1 and 10"}`, http.StatusBadRequest)
		return
	}

	params := models.CreateReviewParams{
		UserID:           userID,
		MovieID:          movieID,
		Rating:           request.Rating,
		Title:            request.Title,
		Content:          request.Content,
		ContainsSpoilers: false,
		IsPublic:         true,
	}

	review, err := h.reviewRepo.CreateReview(r.Context(), params)
	if err != nil {
		log.Printf("Error creating review via API: %v", err)
		http.Error(w, `{"error": "Failed to create review"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
