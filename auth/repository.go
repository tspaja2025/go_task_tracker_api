package auth

import (
	"context"
	"fmt"
	"main/models"

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
