package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "product created successfully"}`))
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"products": [{"id": 1, "name": "dummy product"}]}`))
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	response := fmt.Sprintf(`{"id": "%s", "name": "dummy product %s", "quantity": 10}`, id, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	response := fmt.Sprintf(`{"message": "product %s updated successfully"}`, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	response := fmt.Sprintf(`{"message": "product %s deleted successfully"}`, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func IncrementProductStock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	response := fmt.Sprintf(`{"message": "stock for product %s incremented"}`, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
