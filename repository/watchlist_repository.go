package repository

import (
	"battleNet/models"
	"context"

	"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistRepository struct {
	pool *pgxpool.Pool
}

func NewWatchlistRepository(pool *pgxpool.Pool) *WatchlistRepository {
	return &WatchlistRepository{pool: pool}
}

// AddToWatchlist adds a movie to user's watchlist
func (r *WatchlistRepository) AddToWatchlist(ctx context.Context, userID, movieID uuid.UUID) (*models.WatchlistItem, error) {
	query := `
		INSERT INTO watch_list (user_id, movie_id)
		VALUES ($1, $2)
		RETURNING watch_list_id, user_id, movie_id, added_at
	`

	var item models.WatchlistItem
	err := r.pool.QueryRow(ctx, query, userID, movieID).Scan(
		&item.WatchListID, &item.UserID, &item.MovieID, &item.AddedAt,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetUserWatchlist gets user's watchlist with movie details
func (r *WatchlistRepository) GetUserWatchlist(ctx context.Context, userID uuid.UUID) ([]models.WatchlistItem, error) {
	query := `
		SELECT m.movie_id, m.imdb_id, m.title, m.overview, m.release_date,
			   m.poster_path, m.backdrop_path, m.vote_average, m.vote_count,
			   m.popularity, m.runtime, m.status, m.created_at, w.added_at
		FROM watch_list w
		JOIN movie m ON w.movie_id = m.movie_id
		WHERE w.user_id = $1
		ORDER BY w.added_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.WatchlistItem
	for rows.Next() {
		var item models.WatchlistItem
		var movie models.Movie

		err := rows.Scan(
			&movie.MovieID, &movie.ImdbID, &movie.Title, &movie.Overview, &movie.ReleaseDate,
			&movie.PosterPath, &movie.BackdropPath, &movie.VoteAverage, &movie.VoteCount,
			&movie.Popularity, &movie.Runtime, &movie.Status, &movie.CreatedAt,
			&item.AddedAt,
		)
		if err != nil {
			return nil, err
		}

		item.Movie = movie
		item.UserID = userID
		item.MovieID = movie.MovieID
		items = append(items, item)
	}

	return items, nil
}

// RemoveFromWatchlist removes a movie from user's watchlist
func (r *WatchlistRepository) RemoveFromWatchlist(ctx context.Context, userID, movieID uuid.UUID) error {
	query := `DELETE FROM watch_list WHERE user_id = $1 AND movie_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, movieID)
	return err
}

// CheckWatchlist checks if a movie is in user's watchlist
func (r *WatchlistRepository) CheckWatchlist(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM watch_list WHERE user_id = $1 AND movie_id = $2`

	var count int
	err := r.pool.QueryRow(ctx, query, userID, movieID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
