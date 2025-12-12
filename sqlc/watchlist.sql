-- name: AddToWatchlist :one
INSERT INTO watch_list (user_id, movie_id)
VALUES ($1, $2)
    RETURNING watch_list_id, user_id, movie_id, added_at;

-- name: GetUserWatchlist :many
SELECT m.movie_id, m.imdb_id, m.title, m.overview, m.release_date,
       m.poster_path, m.backdrop_path, m.vote_average, m.vote_count,
       m.popularity, m.runtime, m.status, m.created_at, w.added_at
FROM watch_list w
         JOIN movie m ON w.movie_id = m.movie_id
WHERE w.user_id = $1
ORDER BY w.added_at DESC;

-- name: RemoveFromWatchlist :exec
DELETE FROM watch_list
WHERE user_id = $1 AND movie_id = $2;

-- name: CheckWatchlist :one
SELECT COUNT(*) FROM watch_list
WHERE user_id = $1 AND movie_id = $2;