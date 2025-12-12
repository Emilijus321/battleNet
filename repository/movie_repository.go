package repository

import (
	"battleNet/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{pool: pool}
}

func (r *MovieRepository) GetMovies(ctx context.Context, limit, offset int32) ([]models.Movie, error) {
	query := `
		SELECT movie_id, imdb_id, title, overview, release_date, poster_path,
		       backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at
		FROM movie
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.MovieID, &movie.ImdbID, &movie.Title, &movie.Overview, &movie.ReleaseDate,
			&movie.PosterPath, &movie.BackdropPath, &movie.VoteAverage, &movie.VoteCount,
			&movie.Popularity, &movie.Runtime, &movie.Status, &movie.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *MovieRepository) GetMovieByID(ctx context.Context, movieID uuid.UUID) (*models.Movie, error) {
	query := `
		SELECT movie_id, imdb_id, title, overview, release_date, poster_path,
		       backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at
		FROM movie 
		WHERE movie_id = $1
	`

	var movie models.Movie
	err := r.pool.QueryRow(ctx, query, movieID).Scan(
		&movie.MovieID, &movie.ImdbID, &movie.Title, &movie.Overview, &movie.ReleaseDate,
		&movie.PosterPath, &movie.BackdropPath, &movie.VoteAverage, &movie.VoteCount,
		&movie.Popularity, &movie.Runtime, &movie.Status, &movie.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) CreateMovie(ctx context.Context, movie *models.Movie) error {
	query := `
		INSERT INTO movie (title, overview, release_date, poster_path, backdrop_path,
		                   vote_average, vote_count, popularity, runtime, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING movie_id, created_at
	`

	return r.pool.QueryRow(ctx, query,
		movie.Title, movie.Overview, movie.ReleaseDate, movie.PosterPath, movie.BackdropPath,
		movie.VoteAverage, movie.VoteCount, movie.Popularity, movie.Runtime, movie.Status,
	).Scan(&movie.MovieID, &movie.CreatedAt)
}
