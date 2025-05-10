package main

import (
	"net/http"

	"github.com/dr2cc/URLsShortener.git/internal/handlers"
	"github.com/dr2cc/URLsShortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()

	storageInstance := storage.NewStorage()

	// // Обработчик net/http
	//mux.HandleFunc("POST /{$}", handlers.PostHandler(storageInstance))
	//mux.HandleFunc("GET /{id}", handlers.GetHandler(storageInstance))

	// Обработчик chi
	mux.Post("/", handlers.PostHandler(storageInstance))
	mux.Get("/{id}", handlers.GetHandler(storageInstance))

	http.ListenAndServe(":8080", mux)
}
