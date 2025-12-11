package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"inventory-api/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User

	// 1. Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 2. Validasi Email dan Pass min 6 karakter
	if err := validator.New().Struct(user); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// 3. HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	// 4. Simpan ke Database
	if err := h.Repo.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, "Failed register user", http.StatusInternalServerError)
		return
	}

	// 5. Response Sukses
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}
