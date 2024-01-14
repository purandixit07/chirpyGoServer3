package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/purandixit07/chirpyServer/internal/database"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusOK, struct{}{})
	}
	_, err = cfg.DB.UpdateUserChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}
	respondWithJson(w, http.StatusOK, struct{}{})
}
