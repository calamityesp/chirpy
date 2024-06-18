package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUsersRefresh(w http.ResponseWriter, r *http.Request) {

	refreshtoken := r.Header.Get("Authorization")
	if refreshtoken == "" {
		respondWithError(w, http.StatusUnauthorized, "No refreshtoken given")
		return
	}

	splitAuth := strings.Split(refreshtoken, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "malformed authorization header")
		return
	}

	refreshtoken = splitAuth[1]

	user, err := cfg.DB.GetUserByRefreshToken(refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token doesn't exist")
		return
	}

	cfg.testLog("handler user time: ", user.Refresh_token_expire_time.String())

	// check expiration time of refresh token
	isExpired := cfg.ValidateRefreshToken(&user)
	if !isExpired {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired or revoked")
		return
	}

	//retreive a new access token
	cfg.GetNewJWT(&user)
	respondWithJSON(w, http.StatusOK, user)
}
