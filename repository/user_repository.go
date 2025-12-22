package repository

import (
	"battleNet/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT user_id, email, password_hash, first_name, last_name, username, 
		       role, is_active, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM "user" 
		WHERE email = $1 AND is_active = true
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Username, &user.Role, &user.IsActive, &user.AvatarURL, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO "user" (email, password_hash, first_name, last_name, username, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, created_at, updated_at
	`

	return r.pool.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Username, user.Role,
	).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE "user" SET last_login_at = NOW() WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT user_id, email, password_hash, first_name, last_name, username, 
		       role, is_active, avatar_url, email_verified, created_at, updated_at, last_login_at
		FROM "user" 
		WHERE user_id = $1 AND is_active = true
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Username, &user.Role, &user.IsActive, &user.AvatarURL, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID uuid.UUID, firstName, lastName, username string) error {
	query := `
        UPDATE "user" 
        SET first_name = $2, last_name = $3, username = $4, updated_at = NOW()
        WHERE user_id = $1
    `
	_, err := r.pool.Exec(ctx, query, userID, firstName, lastName, username)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	query := `UPDATE "user" SET password_hash = $2, updated_at = NOW() WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID, newPasswordHash)
	return err
}

// Moderator
func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int32) ([]models.User, error) {
	query := `
        SELECT user_id, email, first_name, last_name, username, 
               role, is_active, avatar_url, email_verified, 
               created_at, updated_at, last_login_at
        FROM "user"
        WHERE is_active = true
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID, &user.Email, &user.FirstName, &user.LastName,
			&user.Username, &user.Role, &user.IsActive, &user.AvatarURL, &user.EmailVerified,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// 🆕 Update user role
func (r *UserRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, newRole string) error {
	query := `UPDATE "user" SET role = $2, updated_at = NOW() WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID, newRole)
	return err
}

// 🆕 Deactivate user (soft delete)
func (r *UserRepository) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE "user" SET is_active = false, updated_at = NOW() WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
