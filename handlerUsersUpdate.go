package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/calamityesp/chirpy/common"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	user := common.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body")
		return
	}

	tokenString := r.Header.Get("Authorization")
	if !strings.HasPrefix(tokenString, "Bearer") {
		respondWithError(w, http.StatusInternalServerError, "Token given was not a Bearer token \n")
	}

	secret := cfg.secret_Key
	claim := &CustomClaims{}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("You fucked up %s\n", err)
	}

	log.Printf("Something Something %s", token.Method.Alg())
}
