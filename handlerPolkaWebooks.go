package main

import (
	"encoding/json"
	"net/http"
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
