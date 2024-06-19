package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolkaUserWebhooks(w http.ResponseWriter, r *http.Request) {
	// check for  user upgrade event
	type Data struct {
		UserId int `json:"user_id"`
	}

	type Hook struct {
		Event string `json:"event"`
		Data  Data
	}

	// check for api key, if not present reject
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "no auth header included in request")
		return
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "malformed authorization header")
		return
	}
	tokenString := splitAuth[1]
	isAuthorized := cfg.CheckApiKey(tokenString)
	if !isAuthorized {
		respondWithError(w, http.StatusUnauthorized, "unauthorized event access")
		return
	}

	params := Hook{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body of requiest")
		return
	}

	switch params.Event {
	case "user.upgraded":
		err = cfg.DB.UpgradeUserToChirpyRed(params.Data.UserId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		break
	case "user.payment_failed":
		err = cfg.DB.DownGradeUserFromChirpyRed(params.Data.UserId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		break
	default:
		respondWithError(w, http.StatusUnauthorized, "Error, invalid event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) CheckApiKey(token string) bool {
	if token == cfg.API_KEY {
		return true
	}
	return false
}
