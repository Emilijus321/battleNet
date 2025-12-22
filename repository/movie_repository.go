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
        INSERT INTO movie (imdb_id, title, overview, release_date, poster_path, backdrop_path,
                           vote_average, vote_count, popularity, runtime, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING movie_id, created_at
    `

	// Jei IMDB ID tuščias, įrašome NULL
	var imdbID interface{}
	if movie.ImdbID != nil && *movie.ImdbID != "" {
		imdbID = *movie.ImdbID
	} else {
		imdbID = nil
	}

	return r.pool.QueryRow(ctx, query,
		imdbID,             // $1 - IMDB ID arba NULL
		movie.Title,        // $2
		movie.Overview,     // $3
		movie.ReleaseDate,  // $4
		movie.PosterPath,   // $5 - gali būti NULL
		movie.BackdropPath, // $6 - gali būti NULL
		movie.VoteAverage,  // $7
		movie.VoteCount,    // $8
		movie.Popularity,   // $9
		movie.Runtime,      // $10
		movie.Status,       // $11
	).Scan(&movie.MovieID, &movie.CreatedAt)
}

// UpdateMovie atnaujina filmą
func (r *MovieRepository) UpdateMovie(ctx context.Context, movieID uuid.UUID, movie *models.Movie) error {
	query := `
        UPDATE movie
        SET title = $2, overview = $3, release_date = $4, 
            vote_average = $5, runtime = $6, status = $7
        WHERE movie_id = $1
    `

	_, err := r.pool.Exec(ctx, query,
		movieID,
		movie.Title,
		movie.Overview,
		movie.ReleaseDate,
		movie.VoteAverage,
		movie.Runtime,
		movie.Status,
	)

	return err
}

// DeleteMovie ištrina filmą
func (r *MovieRepository) DeleteMovie(ctx context.Context, movieID uuid.UUID) error {
	query := `DELETE FROM movie WHERE movie_id = $1`
	_, err := r.pool.Exec(ctx, query, movieID)
	return err
}
