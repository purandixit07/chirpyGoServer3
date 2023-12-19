package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	servmux := http.NewServeMux()
	servmux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	corsMux := middlewareCors(servmux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving the file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
