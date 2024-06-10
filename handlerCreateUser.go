package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	log.Output(1, params.Email)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
	}

	email, err := validateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Email is not valid")
	}

	user, err := cfg.DB.CreateUser(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Id:    user.Id,
		Email: user.Email,
	})
}

func validateEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("No email to validate")
	}

	if !strings.Contains(email, "@") {
		return "", errors.New("Input is not an email address")
	}

	return email, nil

}
