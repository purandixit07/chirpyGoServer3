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

	router := chi.NewRouter()
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", cfg.handlerMetrics)
	apiRouter.Get("/reset", cfg.handlerReset)

	router.Mount("/api", apiRouter)
	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving the file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
