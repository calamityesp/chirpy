package main

import (
	"encoding/json"
	"github.com/calamityesp/chirpy/common"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// get decoder for request body
	decoder := json.NewDecoder(r.Body)
	params := common.User{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode body")
		return
	}

	// get the user from database using the user id
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Not able to find User")
		return
	}

	log.Printf("User password: %s ----  Request password: %s\n", user.Password, params.Password)

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		log.Printf("Comparing hash Error: %s", err.Error())
		log.Printf("req password: %s ----  user.password: %s", params.Password, user.Password)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user.Expires_in_seconds = params.Expires_in_seconds

	// issue a jwt
	cfg.GetNewJWT(&user)

	respondWithJSON(w, http.StatusOK, user)
}
