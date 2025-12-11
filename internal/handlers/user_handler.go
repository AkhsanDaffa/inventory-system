package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 1. Decode & Validasi Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Cari User di Database by Email
	user, err := h.Repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// PENTING: Jangan bilang "Email tidak ditemukan" demi keamanan.
		// Bilang saja "Invalid email or password" agar hacker bingung.
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Cek Password (Bandingkan Hash DB vs Input User)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 4. BIKIN TIKET (JWT) 🎫
	// Menentukan masa berlaku token (misal 24 jam)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Membuat Claims (Isi data di dalam token)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(), // Expired kapan
	}

	// Tanda tangani token dengan JWT_SECRET dari .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		http.Error(w, "Gagal generate token", http.StatusInternalServerError)
		return
	}

	// 5. Kirim Token ke User
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: tokenString,
	})
}
