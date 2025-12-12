package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"inventory-api/internal/database"
	"inventory-api/internal/handlers"
	appMiddleware "inventory-api/internal/middleware"
	"inventory-api/internal/repository"
)

func main() {
	// 1. Setup Structured Logger (JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env Not Found, Use Environment System")
	}

	connString := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	if connString == "" || port == "" {
		slog.Error("Config Environment Incomplete")
		os.Exit(1)
	}

	dbPool, err := database.InitDB(connString)
	if err != nil {
		slog.Error("Could not initialize database: %v", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	productRepo := &repository.ProductRepository{
		DB: dbPool,
	}

	productHandler := &handlers.ProductHandler{
		Repo: productRepo,
	}

	categoryRepo := &repository.CategoryRepository{
		DB: dbPool,
	}

	categoryHandler := &handlers.CategoryHandler{
		Repo: categoryRepo,
	}

	customerRepo := &repository.CustomerRepository{
		DB: dbPool,
	}

	customerHandler := &handlers.CustomerHandler{
		Repo: customerRepo,
	}

	userRepo := &repository.UserRepository{
		DB: dbPool,
	}

	userHandler := &handlers.UserHandler{
		Repo: userRepo,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAllProducts)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware)
			r.Post("/", productHandler.CreateProduct)
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetProductByID)

			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.AuthMiddleware)

				r.Put("/", productHandler.UpdateProduct)
				r.Delete("/", productHandler.DeleteProduct)
			})
		})
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", categoryHandler.GetAllCategories)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware)
			r.Post("/", categoryHandler.CreateCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
		})
	})

	r.Route("/customers", func(r chi.Router) {
		r.Get("/", customerHandler.GetAllCustomers)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.AuthMiddleware)
			r.Post("/", customerHandler.CreateCustomer)
		})
	})

	r.Post("/register", userHandler.RegisterUser)
	r.Post("/login", userHandler.LoginUser)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Inventory API"))
	})

	slog.Info("Starting starting...", "port", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
