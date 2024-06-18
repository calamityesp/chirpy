package main

import (
	"encoding/json"
	"errors"
	"github.com/calamityesp/chirpy/common"
	"net/http"
	"strings"
)

// type User struct {
// 	Id       int    `json:"id"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	params := common.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
	}

	err = validateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Email is not valid")
	}

	test, err := cfg.DB.GetUserByEmail(params.Email)
	if test.Email == params.Email {
		respondWithError(w, http.StatusUnauthorized, "User Already Exists, or error with internal database")
		return
	}

	user, err := cfg.DB.CreateUser(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, common.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	})
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("No email to validate")
	}

	if !strings.Contains(email, "@") {
		return errors.New("Input is not an email address")
	}
	return nil
}
