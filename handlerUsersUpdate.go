package main

import (
	"encoding/json"
	"github.com/calamityesp/chirpy/common"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	user := common.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body")
		return
	}

	// tokenString := r.Header.Get("Authorization")
	// if !strings.HasPrefix(tokenString, "Bearer") {
	// 	respondWithError(w, http.StatusInternalServerError, "Token given was not a Bearer token \n")
	// }

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "no auth header included in request")
		return
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "malformed authorization header")
		return
	}
	tokenString := splitAuth[1]

	secret := cfg.secret_Key
	claim := jwt.RegisteredClaims{}
	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	log.Printf("GivenToken:   %s", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("You fucked up %s\n", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// validate token
	if !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	log.Printf("ExpiresAt: %v", claim.ExpiresAt.UTC())
	// Check if the token is expired
	if claim.ExpiresAt != nil && claim.ExpiresAt.Time.Before(time.Now()) {
		log.Printf("ExpiresAt: %v", claim.ExpiresAt.UTC())
		respondWithError(w, http.StatusUnauthorized, "Token has expired")
		return
	}

	subject, _ := token.Claims.GetSubject()
	user.Id, err = strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to convert id to integer")
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if "chirpy" != issuer {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token.")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Validating Token")
		return
	}
	update, err := cfg.DB.UpdateUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	update.Token = tokenString
	respondWithJSON(w, http.StatusOK, update)

}

// func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Password string `json:"password"`
// 		Email    string `json:"email"`
// 	}
// 	type response struct {
// 		User
// 	}
//
// 	token, err := auth.GetBearerToken(r.Header)
// 	if err != nil {
// 		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
// 		return
// 	}
// 	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
// 	if err != nil {
// 		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
// 		return
// 	}
//
// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err = decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
// 		return
// 	}
//
// 	hashedPassword, err := auth.HashPassword(params.Password)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
// 		return
// 	}
//
// 	userIDInt, err := strconv.Atoi(subject)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
// 		return
// 	}
//
// 	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
// 		return
// 	}
//
// 	respondWithJSON(w, http.StatusOK, response{
// 		User: User{
// 			ID:    user.ID,
// 			Email: user.Email,
// 		},
// 	})
// }
