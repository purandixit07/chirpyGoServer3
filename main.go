package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/purandixit07/chirpyServer/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", cfg.handlerReset)
	apiRouter.Post("/chirps", cfg.handlerCreateChirps)
	apiRouter.Post("/users", cfg.handlerUsersCreate)
	apiRouter.Get("/chirps", cfg.handlerRetrieveChirps)
	apiRouter.Get("/chirps/{chirpID}", cfg.handlerRetrieveChirpById)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.handlerMetrics)

	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)
	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving the file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
