package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"inventory-api/internal/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		r.Post("/", handlers.CreateProduct)
		r.Get("/", handlers.GetAllProducts)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetProductByID)
			r.Put("/", handlers.UpdateProduct)
			r.Delete("/", handlers.DeleteProduct)
			r.Post("/increment-stock", handlers.IncrementProductStock)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Inventory API"))
	})

	log.Println("Starting server on :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
