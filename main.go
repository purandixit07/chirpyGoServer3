package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/purandixit07/chirpyServer/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	cfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	router := chi.NewRouter()
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", cfg.handlerReset)

	apiRouter.Post("/login", cfg.handlerUserLogin)
	apiRouter.Post("/refresh", cfg.handlerRefresh)
	apiRouter.Post("/revoke", cfg.handlerRevoke)

	apiRouter.Post("/users", cfg.handlerUsersCreate)
	apiRouter.Put("/users", cfg.handlerUsersUpdate)

	apiRouter.Post("/chirps", cfg.handlerCreateChirps)
	apiRouter.Get("/chirps", cfg.handlerRetrieveChirps)
	apiRouter.Get("/chirps/{chirpID}", cfg.handlerRetrieveChirpById)
	apiRouter.Delete("/chirps/{chirpID}", cfg.handlerDeleteChirp)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving the file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
