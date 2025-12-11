package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) CreateUser(ctx context.Context, u *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	err := r.DB.QueryRow(ctx, query, u.Email, u.Password).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("failed register user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	var u User
	err := r.DB.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &u, nil
}
