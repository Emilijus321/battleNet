-- +goose Up
-- +goose StatementBegin
ALTER TABLE movie_genre
    ADD CONSTRAINT fk_movie_genre_movie FOREIGN KEY (movie_id) REFERENCES movie(movie_id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_movie_genre_genre FOREIGN KEY (genre_id) REFERENCES genre(genre_id) ON DELETE CASCADE;

ALTER TABLE oauth
    ADD CONSTRAINT fk_oauth_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE;

ALTER TABLE review
    ADD CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_review_movie FOREIGN KEY (movie_id) REFERENCES movie(movie_id) ON DELETE CASCADE;

ALTER TABLE review_like
    ADD CONSTRAINT fk_review_like_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_review_like_review FOREIGN KEY (review_id) REFERENCES review(review_id) ON DELETE CASCADE;

ALTER TABLE user_session
    ADD CONSTRAINT fk_user_session_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE;

ALTER TABLE watch_list
    ADD CONSTRAINT fk_watch_list_user FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_watch_list_movie FOREIGN KEY (movie_id) REFERENCES movie(movie_id) ON DELETE CASCADE;

CREATE INDEX idx_user_email ON "user"(email);
CREATE INDEX idx_user_username ON "user"(username);
CREATE INDEX idx_movie_imdb_id ON movie(imdb_id);
CREATE INDEX idx_movie_title ON movie(title);
CREATE INDEX idx_movie_release_date ON movie(release_date);
CREATE INDEX idx_review_user_id ON review(user_id);
CREATE INDEX idx_review_movie_id ON review(movie_id);
CREATE INDEX idx_review_created_at ON review(created_at);
CREATE INDEX idx_watch_list_user_id ON watch_list(user_id);
CREATE INDEX idx_watch_list_movie_id ON watch_list(movie_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_watch_list_movie_id;
DROP INDEX IF EXISTS idx_watch_list_user_id;
DROP INDEX IF EXISTS idx_review_created_at;
DROP INDEX IF EXISTS idx_review_movie_id;
DROP INDEX IF EXISTS idx_review_user_id;
DROP INDEX IF EXISTS idx_movie_release_date;
DROP INDEX IF EXISTS idx_movie_title;
DROP INDEX IF EXISTS idx_movie_imdb_id;
DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_email;

ALTER TABLE watch_list DROP CONSTRAINT IF EXISTS fk_watch_list_movie;
ALTER TABLE watch_list DROP CONSTRAINT IF EXISTS fk_watch_list_user;
ALTER TABLE user_session DROP CONSTRAINT IF EXISTS fk_user_session_user;
ALTER TABLE review_like DROP CONSTRAINT IF EXISTS fk_review_like_review;
ALTER TABLE review_like DROP CONSTRAINT IF EXISTS fk_review_like_user;
ALTER TABLE review DROP CONSTRAINT IF EXISTS fk_review_movie;
ALTER TABLE review DROP CONSTRAINT IF EXISTS fk_review_user;
ALTER TABLE oauth DROP CONSTRAINT IF EXISTS fk_oauth_user;
ALTER TABLE movie_genre DROP CONSTRAINT IF EXISTS fk_movie_genre_genre;
ALTER TABLE movie_genre DROP CONSTRAINT IF EXISTS fk_movie_genre_movie;
-- +goose StatementEnd