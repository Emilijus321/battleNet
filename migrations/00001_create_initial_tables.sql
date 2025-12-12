-- migrations/00001_create_initial_tables.sql
-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users table
CREATE TABLE "user" (
                        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        email VARCHAR(255) UNIQUE NOT NULL,
                        password_hash TEXT NOT NULL,
                        first_name VARCHAR(100) NOT NULL,
                        last_name VARCHAR(100) NOT NULL,
                        username VARCHAR(100) UNIQUE NOT NULL,
                        role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
                        is_active BOOLEAN DEFAULT true,
                        avatar_url TEXT,
                        email_verified BOOLEAN DEFAULT false,
                        created_at TIMESTAMPTZ DEFAULT NOW(),
                        updated_at TIMESTAMPTZ DEFAULT NOW(),
                        last_login_at TIMESTAMPTZ
);

-- 2. Genre table
CREATE TABLE genre (
                       genre_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name VARCHAR(100) UNIQUE NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Movie table
CREATE TABLE movie (
                       movie_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       imdb_id VARCHAR(20) UNIQUE,
                       title VARCHAR(500) NOT NULL,
                       overview TEXT,
                       release_date DATE,
                       poster_path TEXT,
                       backdrop_path TEXT,
                       vote_average DECIMAL(3,1) DEFAULT 0.0,
                       vote_count INTEGER DEFAULT 0,
                       popularity DECIMAL(10,4) DEFAULT 0.0,
                       runtime INTEGER, -- in minutes
                       status VARCHAR(50) DEFAULT 'Released', -- Released, Post Production, In Production, etc.
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Movie-Genre relationship (many-to-many)
CREATE TABLE movie_genre (
                             movie_genre_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             movie_id UUID NOT NULL,
                             genre_id UUID NOT NULL,
                             created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. OAuth providers table
CREATE TABLE oauth (
                       oauth_provider_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL,
                       provider VARCHAR(50) NOT NULL, -- 'google', 'facebook', 'github', etc.
                       provider_id VARCHAR(255) NOT NULL, -- ID from the OAuth provider
                       provider_email VARCHAR(255),
                       access_token TEXT,
                       refresh_token TEXT,
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Reviews table
CREATE TABLE review (
                        review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NOT NULL,
                        movie_id UUID NOT NULL,
                        rating INTEGER CHECK (rating >= 1 AND rating <= 10),
                        title VARCHAR(255) NOT NULL,
                        content TEXT NOT NULL,
                        contains_spoilers BOOLEAN DEFAULT false,
                        is_public BOOLEAN DEFAULT true,
                        likes_count INTEGER DEFAULT 0,
                        created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Review likes table
CREATE TABLE review_like (
                             review_like_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             user_id UUID NOT NULL,
                             review_id UUID NOT NULL,
                             is_like BOOLEAN DEFAULT true, -- true for like, false for dislike
                             created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. User sessions table
CREATE TABLE user_session (
                              session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              user_id UUID NOT NULL,
                              session_token VARCHAR(255) UNIQUE NOT NULL,
                              device_info TEXT,
                              ip_address INET,
                              expires_at TIMESTAMPTZ NOT NULL,
                              created_at TIMESTAMPTZ DEFAULT NOW(),
                              last_used_at TIMESTAMPTZ DEFAULT NOW()
);

-- 9. Watchlist table
CREATE TABLE watch_list (
                            watch_list_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            user_id UUID NOT NULL,
                            movie_id UUID NOT NULL,
                            added_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS watch_list;
DROP TABLE IF EXISTS user_session;
DROP TABLE IF EXISTS review_like;
DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS oauth;
DROP TABLE IF EXISTS movie_genre;
DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS "user";

-- +goose StatementEnd