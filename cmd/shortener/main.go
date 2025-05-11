package main

import (
	"net/http"

	"github.com/dr2cc/URLsShortener.git/internal/handlers"
	"github.com/dr2cc/URLsShortener.git/internal/storage"
)

func main() {
	mux := http.NewServeMux()
	storageInstance := storage.NewStorage()

	mux.HandleFunc("POST /{$}", handlers.PostHandler(storageInstance))
	//mux.HandleFunc("GET /{id}", handlers.GetHandler(storageInstance))

	// Работает и так и так.
	// Автотесты прошел! Верхний конечно логичнее.
	// Пока оставлю и так и так, как упражнение
	mux.HandleFunc("GET /{id}", storageInstance.GetHandler)

	http.ListenAndServe(":8080", mux)
}
