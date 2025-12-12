-- name: CreateReview :one
INSERT INTO review (user_id, movie_id, rating, title, content, contains_spoilers, is_public)
VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING review_id, user_id, movie_id, rating, title, content,
          contains_spoilers, is_public, likes_count, created_at;

-- name: GetReviewByID :one
SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
       r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
       u.username, u.avatar_url
FROM review r
         JOIN "user" u ON r.user_id = u.user_id
WHERE r.review_id = $1;

-- name: GetMovieReviews :many
SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
       r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
       u.username, u.avatar_url
FROM review r
         JOIN "user" u ON r.user_id = u.user_id
WHERE r.movie_id = $1 AND r.is_public = true
ORDER BY r.created_at DESC;

-- name: GetUserReviews :many
SELECT r.review_id, r.user_id, r.movie_id, r.rating, r.title, r.content,
       r.contains_spoilers, r.is_public, r.likes_count, r.created_at,
       m.title as movie_title, m.poster_path
FROM review r
         JOIN movie m ON r.movie_id = m.movie_id
WHERE r.user_id = $1
ORDER BY r.created_at DESC;

-- name: UpdateReview :one
UPDATE review
SET rating = $2, title = $3, content = $4, contains_spoilers = $5, is_public = $6
WHERE review_id = $1
    RETURNING review_id, user_id, movie_id, rating, title, content,
          contains_spoilers, is_public, likes_count, created_at;

-- name: DeleteReview :exec
DELETE FROM review WHERE review_id = $1;