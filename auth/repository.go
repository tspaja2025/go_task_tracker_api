package auth

import (
	"context"
	"fmt"
	"main/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Hash the password and insert the new record into Postgres
func (r *Repository) CreateUser(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	// Hash the password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(hashedBytes)

	// Insert into database
	query := `
					INSERT INTO users (username, email, password_hash)
					VALUES ($1, $2, $3)
					RETURNING id, username, email, created_at, updated_at;
	`

	var user models.User
	err = r.db.QueryRow(ctx, query, req.Username, req.Email, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Protect againts duplicate emails or usernames
		return nil, err
	}

	return &user, nil
}

// Search for user by email to verify identity during login
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1;`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		// Will return pgx.ErrNoRows if user does not exits
		return nil, err
	}
	return &user, nil
}

// Store a long-lived session token hash in the database
func (r *Repository) SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}
