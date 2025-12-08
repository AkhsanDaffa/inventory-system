package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name" validate:"required"`
	SKU        string `json:"sku" validate:"required"`
	Quantity   int    `json:"quantity" validate:"gte=0"`
	CategoryID string `json:"category_id"`

	CategoryName string `json:"category_name,omitempty"`
}

type ProductRepository struct {
	DB *pgxpool.Pool
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *Product) error {
	query := `
		INSERT INTO products (name, sku, quantity, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.DB.QueryRow(ctx, query, p.Name, p.SKU, p.Quantity, p.CategoryID).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("failed Insert Database: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]Product, error) {
	products := []Product{}

	// query := "SELECT id, name, sku, quantity FROM products"
	query := `
	SELECT
		p.id, p.name, p.sku, p.quantity, p.category_id,
		COALESCE(c.name, '') as category_name
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		// if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Quantity); err != nil {
		// 	return nil, fmt.Errorf("failed to scan: %w", err)
		// }
		var catID *string
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Quantity, &catID, &p.CategoryName); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		if catID != nil {
			p.CategoryID = *catID
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (Product, error) {
	var p Product
	query := "SELECT id, name, sku, quantity FROM products WHERE id = $1"

	err := r.DB.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.SKU, &p.Quantity)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, p *Product) error {
	query := "UPDATE products SET name=$1, sku=$2, quantity=$3 WHERE id=$4"

	commandTag, err := r.DB.Exec(ctx, query, p.Name, p.SKU, p.Quantity, id)
	if err != nil {
		return fmt.Errorf("failed Update: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id=$1"

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
