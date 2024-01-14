package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	authorID := -1
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author id")
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorID != -1 && dbChirp.AuthorID != authorID {
			continue
		}
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerRetrieveChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID := chi.URLParam(r, "chirpID")
	chirpIDInt, err := strconv.Atoi(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID")
		return
	}
	dbChirp, err := cfg.DB.GetChirpByID(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})

}
