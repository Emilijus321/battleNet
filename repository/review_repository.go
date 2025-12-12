package repository

import (
	"battleNet/models"
	"context"

	"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{pool: pool}
}

// CreateReview creates a new review
func (r *ReviewRepository) CreateReview(ctx context.Context, params models.CreateReviewParams) (*models.Review, error) {
	query := `
		INSERT INTO review (user_id, movie_id, rating, title, content, contains_spoilers, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING review_id, user_id, movie_id, rating, title, content,
				  contains_spoilers, is_public, likes_count, created_at
	`

	var review models.Review
	err := r.pool.QueryRow(ctx, query,
		params.UserID, params.MovieID, params.Rating, params.Title, params.Content,
		params.ContainsSpoilers, params.IsPublic,
	).Scan(
		&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Title,
		&review.Content, &review.ContainsSpoilers, &review.IsPublic, &review.LikesCount,
		&review.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetReviewByID gets a review by ID with user information
func (r *ReviewRepository) GetReviewByID(ctx context.Context, reviewID uuid.UUID) (*models.Review, error) {
	query := `
		SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
			   r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
			   u.username, u.avatar_url
		FROM review r
		JOIN "user" u ON r.user_id = u.user_id
		WHERE r.review_id = $1
	`

	var review models.Review
	err := r.pool.QueryRow(ctx, query, reviewID).Scan(
		&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Title,
		&review.Content, &review.ContainsSpoilers, &review.IsPublic, &review.LikesCount,
		&review.CreatedAt, &review.Username, &review.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetMovieReviews gets all public reviews for a movie
func (r *ReviewRepository) GetMovieReviews(ctx context.Context, movieID uuid.UUID) ([]models.Review, error) {
	query := `
		SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
			   r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
			   u.username, u.avatar_url
		FROM review r
		JOIN "user" u ON r.user_id = u.user_id
		WHERE r.movie_id = $1 AND r.is_public = true
		ORDER BY r.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Title,
			&review.Content, &review.ContainsSpoilers, &review.IsPublic, &review.LikesCount,
			&review.CreatedAt, &review.Username, &review.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// GetUserReviews gets all reviews by a user
func (r *ReviewRepository) GetUserReviews(ctx context.Context, userID uuid.UUID) ([]models.Review, error) {
	query := `
		SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
			   r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
			   m.title as movie_title, m.poster_path
		FROM review r
		JOIN movie m ON r.movie_id = m.movie_id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var movieTitle string
		var posterPath *string

		err := rows.Scan(
			&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Title,
			&review.Content, &review.ContainsSpoilers, &review.IsPublic, &review.LikesCount,
			&review.CreatedAt, &movieTitle, &posterPath,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// UpdateReview updates a review
func (r *ReviewRepository) UpdateReview(ctx context.Context, reviewID uuid.UUID, params models.UpdateReviewParams) (*models.Review, error) {
	query := `
		UPDATE review
		SET rating = $2, title = $3, content = $4, contains_spoilers = $5, is_public = $6
		WHERE review_id = $1
		RETURNING review_id, user_id, movie_id, rating, title, content,
				  contains_spoilers, is_public, likes_count, created_at
	`

	var review models.Review
	err := r.pool.QueryRow(ctx, query,
		reviewID, params.Rating, params.Title, params.Content,
		params.ContainsSpoilers, params.IsPublic,
	).Scan(
		&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Title,
		&review.Content, &review.ContainsSpoilers, &review.IsPublic, &review.LikesCount,
		&review.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

// DeleteReview deletes a review
func (r *ReviewRepository) DeleteReview(ctx context.Context, reviewID uuid.UUID) error {
	query := `DELETE FROM review WHERE review_id = $1`
	_, err := r.pool.Exec(ctx, query, reviewID)
	return err
}
