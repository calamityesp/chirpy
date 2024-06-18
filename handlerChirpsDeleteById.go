package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsDeleteById(w http.ResponseWriter, r *http.Request) {
	// get jwt token
	authHeader := r.Header.Get("Authorization")
	user, isValid := cfg.validateJwtToken(authHeader)
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Invalid authentication token")
		return
	}

	chirpId := r.PathValue("chirpId")
	if chirpId == "" {
		respondWithError(w, http.StatusUnauthorized, "no chirp id present")
		return
	}

	chirpInt, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	chirp, err := cfg.DB.GetChirpById(chirpInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	//compoare chirp with user handlerChirpsDeleteById
	if user.Id != chirp.Author_Id {
		respondWithError(w, http.StatusForbidden, "invalid: user is not author of chirp")
		return
	}

	// delete the handlerChirpsDeleteById
	err = cfg.DB.DeleteChirpById(chirpInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting Chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
