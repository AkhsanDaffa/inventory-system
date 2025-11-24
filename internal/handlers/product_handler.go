package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHandler struct {
	DB *pgxpool.Pool
}

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte(`{"message": "product created successfully"}`))

	var product Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO products (name, sku, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var newID string
	err := h.DB.QueryRow(r.Context(), query, product.Name, product.SKU, product.Quantity).Scan(&newID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed store to DB: %v", err), http.StatusInternalServerError)
		return
	}

	product.ID = newID

	// fmt.Printf("Received product: %+v\n", product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product created successfully",
		"data":    product,
	})
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, name, sku, quantity FROM products")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product

		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Quantity); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan product: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": products,
	})
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// response := fmt.Sprintf(`{"id": "%s", "name": "dummy product %s", "quantity": 10}`, id, id)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(response))

	var p Product

	query := "SELECT id, name, sku, quantity FROM products WHERE id=$1"

	err := h.DB.QueryRow(r.Context(), query, id).Scan(&p.ID, &p.Name, &p.SKU, &p.Quantity)

	if err != nil {
		http.Error(w, "Product not found or invalid ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": p,
	})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// response := fmt.Sprintf(`{"message": "product %s updated successfully"}`, id)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(response))

	var p Product

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if p.Name == "" || p.SKU == "" {
		http.Error(w, "Name and SKU are required fields", http.StatusBadRequest)
		return
	}

	query := "UPDATE products SET name=$1, sku=$2, quantity=$3 WHERE id=$4"

	commandTag, err := h.DB.Exec(r.Context(), query, p.Name, p.SKU, p.Quantity, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update product: %v", err), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product updated successfully",
	})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// response := fmt.Sprintf(`{"message": "product %s deleted successfully"}`, id)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(response))

	query := "DELETE FROM products WHERE id=$1"

	commandTag, err := h.DB.Exec(r.Context(), query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete product: %v", err), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product deleted successfully",
	})
}

// func IncrementProductStock(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")

// 	response := fmt.Sprintf(`{"message": "stock for product %s incremented"}`, id)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(response))
// }
