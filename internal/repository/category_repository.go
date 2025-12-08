package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
}

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c *Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`

	err := r.DB.QueryRow(ctx, query, c.Name).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed insert category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]Category, error) {
	categories := []Category{}

	query := `SELECT id, name FROM categories`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
