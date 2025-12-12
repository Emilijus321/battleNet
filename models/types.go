package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	FirstName     string     `json:"first_name" db:"first_name"`
	LastName      string     `json:"last_name" db:"last_name"`
	Username      string     `json:"username" db:"username"`
	Role          string     `json:"role" db:"role"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	AvatarURL     *string    `json:"avatar_url" db:"avatar_url"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at" db:"last_login_at"`
}

type Movie struct {
	MovieID      uuid.UUID  `json:"movie_id" db:"movie_id"`
	ImdbID       *string    `json:"imdb_id" db:"imdb_id"`
	Title        string     `json:"title" db:"title"`
	Overview     *string    `json:"overview" db:"overview"`
	ReleaseDate  *time.Time `json:"release_date" db:"release_date"`
	PosterPath   *string    `json:"poster_path" db:"poster_path"`
	BackdropPath *string    `json:"backdrop_path" db:"backdrop_path"`
	VoteAverage  *float64   `json:"vote_average" db:"vote_average"`
	VoteCount    *int       `json:"vote_count" db:"vote_count"`
	Popularity   *float64   `json:"popularity" db:"popularity"`
	Runtime      *int       `json:"runtime" db:"runtime"`
	Status       *string    `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

type Review struct {
	ReviewID         uuid.UUID `json:"review_id" db:"review_id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	MovieID          uuid.UUID `json:"movie_id" db:"movie_id"`
	Rating           int       `json:"rating" db:"rating"`
	Title            string    `json:"title" db:"title"`
	Content          string    `json:"content" db:"content"`
	ContainsSpoilers bool      `json:"contains_spoilers" db:"contains_spoilers"`
	IsPublic         bool      `json:"is_public" db:"is_public"`
	LikesCount       int       `json:"likes_count" db:"likes_count"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Username         string    `json:"username" db:"username"`
	AvatarURL        *string   `json:"avatar_url" db:"avatar_url"`
}

type WatchlistItem struct {
	WatchListID uuid.UUID `json:"watch_list_id" db:"watch_list_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	MovieID     uuid.UUID `json:"movie_id" db:"movie_id"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`
	Movie       Movie     `json:"movie" db:"movie"`
}

// Request/Response types
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
}

type CreateReviewRequest struct {
	MovieID string `json:"movie_id"`
	Rating  int    `json:"rating"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Repository parameter types
type CreateReviewParams struct {
	UserID           uuid.UUID
	MovieID          uuid.UUID
	Rating           int
	Title            string
	Content          string
	ContainsSpoilers bool
	IsPublic         bool
}

type UpdateReviewParams struct {
	Rating           int
	Title            string
	Content          string
	ContainsSpoilers bool
	IsPublic         bool
}
