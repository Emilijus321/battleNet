-- +goose Up
-- +goose StatementBegin

-- Insert genres
INSERT INTO genre (name) VALUES
                             ('Action'),
                             ('Adventure'),
                             ('Animation'),
                             ('Comedy'),
                             ('Crime'),
                             ('Documentary'),
                             ('Drama'),
                             ('Family'),
                             ('Fantasy'),
                             ('History'),
                             ('Horror'),
                             ('Music'),
                             ('Mystery'),
                             ('Romance'),
                             ('Science Fiction'),
                             ('TV Movie'),
                             ('Thriller'),
                             ('War'),
                             ('Western');

-- Insert default admin user (password: admin123)
-- Password hash generated with bcrypt
INSERT INTO "user" (
    email,
    password_hash,
    first_name,
    last_name,
    username,
    role,
    email_verified
) VALUES (
             'admin@movieapp.com',
             '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
             'System',
             'Administrator',
             'admin',
             'admin',
             true
         );

-- Insert sample movies
INSERT INTO movie (
    imdb_id,
    title,
    overview,
    release_date,
    poster_path,
    backdrop_path,
    vote_average,
    vote_count,
    popularity,
    runtime,
    status
) VALUES
      (
          'tt1375666',
          'Inception',
          'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.',
          '2010-07-16',
          '/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg',
          '/s2bT29y0ngXxxu2IA8AOzzXTRhd.jpg',
          8.4,
          35000,
          100.5,
          148,
          'Released'
      ),
      (
          'tt0468569',
          'The Dark Knight',
          'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.',
          '2008-07-18',
          '/qJ2tW6WMUDux911r6m7haRef0WH.jpg',
          '/hqkIcbrOHL86UncnHIsHVcVmzue.jpg',
          9.0,
          28000,
          95.2,
          152,
          'Released'
      );

-- Link movies to genres
INSERT INTO movie_genre (movie_id, genre_id)
SELECT m.movie_id, g.genre_id
FROM movie m, genre g
WHERE m.title = 'Inception' AND g.name IN ('Action', 'Science Fiction', 'Thriller');

INSERT INTO movie_genre (movie_id, genre_id)
SELECT m.movie_id, g.genre_id
FROM movie m, genre g
WHERE m.title = 'The Dark Knight' AND g.name IN ('Action', 'Crime', 'Drama');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove sample data
DELETE FROM movie_genre WHERE movie_id IN (
    SELECT movie_id FROM movie WHERE title IN ('Inception', 'The Dark Knight')
);

DELETE FROM movie WHERE title IN ('Inception', 'The Dark Knight');

DELETE FROM "user" WHERE email = 'admin@movieapp.com';

DELETE FROM genre WHERE name IN (
                                 'Action', 'Adventure', 'Animation', 'Comedy', 'Crime', 'Documentary',
                                 'Drama', 'Family', 'Fantasy', 'History', 'Horror', 'Music', 'Mystery',
                                 'Romance', 'Science Fiction', 'TV Movie', 'Thriller', 'War', 'Western'
    );

-- +goose StatementEnd