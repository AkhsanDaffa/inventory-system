package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type CustomerRepository struct {
	DB *pgxpool.Pool
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, c *Customer) error {
	query := `
		INSERT INTO customers (name, email, phone)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.DB.QueryRow(ctx, query, c.Name, c.Email, c.Phone).Scan(&c.ID)

	if err != nil {
		return fmt.Errorf("failed insert database: %w", err)
	}

	return nil
}

func (r *CustomerRepository) GetAllCustomers(ctx context.Context) ([]Customer, error) {
	customers := []Customer{}

	query := "SELECT id, name, email, phone FROM customers"
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		customers = append(customers, c)
	}
	return customers, nil
}
