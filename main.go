package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	cfg := apiConfig{
		fileServerHits: 0,
	}

	r := chi.NewRouter()
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", cfg.handlerMetrics)
	r.Get("/reset", cfg.handlerReset)

	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving the file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
