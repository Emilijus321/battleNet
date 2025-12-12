-- name: GetMovies :many
SELECT movie_id, imdb_id, title, overview, release_date, poster_path,
       backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at
FROM movie
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: GetMoviesCount :one
SELECT COUNT(*) FROM movie;

-- name: GetMovieByID :one
SELECT movie_id, imdb_id, title, overview, release_date, poster_path,
       backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at
FROM movie WHERE movie_id = $1;

-- name: GetMovieByIMDBID :one
SELECT movie_id, imdb_id, title, overview, release_date, poster_path,
       backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at
FROM movie WHERE imdb_id = $1;

-- name: CreateMovie :one
INSERT INTO movie (imdb_id, title, overview, release_date, poster_path,
                   backdrop_path, vote_average, vote_count, popularity, runtime, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING movie_id, imdb_id, title, overview, release_date, poster_path,
          backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at;

-- name: UpdateMovie :one
UPDATE movie
SET imdb_id = $2, title = $3, overview = $4, release_date = $5,
    poster_path = $6, backdrop_path = $7, vote_average = $8,
    vote_count = $9, popularity = $10, runtime = $11, status = $12
WHERE movie_id = $1
    RETURNING movie_id, imdb_id, title, overview, release_date, poster_path,
          backdrop_path, vote_average, vote_count, popularity, runtime, status, created_at;

-- name: DeleteMovie :exec
DELETE FROM movie WHERE movie_id = $1;