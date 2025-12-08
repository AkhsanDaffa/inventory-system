package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"inventory-api/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	Repo *repository.CustomerRepository
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer repository.Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(customer); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateCustomer(r.Context(), &customer)

	if err != nil {
		http.Error(w, "Failed store data (email might be duplicate)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "customer created successfully",
		"data":    customer,
	})
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Repo.GetAllCustomers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": customers,
	})
}
