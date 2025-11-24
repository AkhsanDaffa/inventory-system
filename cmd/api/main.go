package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"inventory-api/internal/database"
	"inventory-api/internal/handlers"
)

func main() {
	connString := "postgres://postgres:testing@localhost:5433/postgres?sslmode=disable"

	dbPool, err := database.InitDB(connString)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	defer dbPool.Close()

	productHandler := &handlers.ProductHandler{
		DB: dbPool,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		// r.Post("/", handlers.CreateProduct)
		// r.Get("/", handlers.GetAllProducts)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetAllProducts)

		// r.Route("/{id}", func(r chi.Router) {
		// 	r.Get("/", handlers.GetProductByID)
		// 	r.Put("/", handlers.UpdateProduct)
		// 	r.Delete("/", handlers.DeleteProduct)
		// 	r.Post("/increment-stock", handlers.IncrementProductStock)
		// })

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetProductByID)
			r.Put("/", productHandler.UpdateProduct)
			r.Delete("/", productHandler.DeleteProduct)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Inventory API"))
	})

	log.Println("Starting server on :8081")

	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
