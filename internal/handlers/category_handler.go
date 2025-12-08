package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"inventory-api/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	Repo *repository.CategoryRepository
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category repository.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(category); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateCategory(r.Context(), &category); err != nil {
		http.Error(w, "Failed save category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category created",
		"data":    category,
	})
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Repo.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, "Faild collect data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": categories,
	})
}
