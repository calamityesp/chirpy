package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUsersRevoke(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		respondWithError(w, http.StatusUnauthorized, "No token given to revoke")
		return
	}

	splitAuth := strings.Split(authToken, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "malformed authorization header")
		return
	}

	// revoke the token
	revoked, err := cfg.DB.RevokeUserRefreshToken(splitAuth[1])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if revoked {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	respondWithError(w, http.StatusUnauthorized, "unknown problem revoking token")

}
