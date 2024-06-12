package main

import (
	"encoding/json"
	"net/http"

	"github.com/calamityesp/chirpy/common"
	"golang.org/x/crypto/bcrypt"
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

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Password")
		return
	}

	// issue a jwt
	cfg.GetNewJWT(&user)

	respondWithJSON(w, http.StatusOK, user)
}
